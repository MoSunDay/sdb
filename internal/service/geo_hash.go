package service

import (
	"errors"
	"fmt"
	"github.com/gansidui/geohash"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/store"
	"google.golang.org/protobuf/proto"
	"math"
	"sort"
	"strconv"
)

var NotFoundGeoHashError = errors.New("not found geo hash, please create it")
var GeoHashExistError = errors.New("geo hash exist, please delete it or change other")

const geoHashMetaKeyTemplate = "gh/%s"

// id key
const geoHashIdPrefixKeyTemplate = "ghi/%s"
const geoHashIdKeyTemplate = geoHashIdPrefixKeyTemplate + "/%s"

// hash key
const geoHashPrefixKeyTemplate = "ghh/%s/%s"
const geoHashKeyTemplate = geoHashPrefixKeyTemplate + "/%s"

// box key
const geoHashBoxPrefixKeyTemplate = "ghb/%s/%s"
const geoHashBoxKeyTemplate = geoHashBoxPrefixKeyTemplate + "/%s"

func GHCreate(key []byte, precision int32) (bool, error) {
	lock(LGeoHash, key)
	defer unlock(LGeoHash, key)

	metaKey := generateGeoHashMetaKey(key)
	exist, err := store.Get(metaKey)
	if err != nil {
		return false, err
	}
	if exist != nil || len(exist) > 0 {
		return false, GeoHashExistError
	}

	return store.Set(metaKey, []byte(strconv.Itoa(int(precision))))
}

func GHDel(key []byte) (bool, error) {
	batch := store.NewBatch()
	defer batch.Close()

	metaKey := generateGeoHashMetaKey(key)
	precisionStr, err := store.Get(metaKey)
	if err != nil {
		return false, err
	}
	if precisionStr == nil || len(precisionStr) == 0 {
		return false, NotFoundGeoHashError
	}
	precision, err := strconv.ParseInt(string(precisionStr), 10, 32)
	if err != nil {
		return false, err
	}

	if _, err := batch.Del(metaKey); err != nil {
		return false, err
	}
	if err := store.Iterate0(generateGeoHashIdPrefixKey(key),
		func(key []byte, value []byte) error {
			_, err := batch.Del(key)
			if err != nil {
				return err
			}

			var existPoint pb.Point
			if err := proto.Unmarshal(value, &existPoint); err != nil {
				return err
			}
			hash, box := geohash.Encode(existPoint.Latitude, existPoint.Longitude, int(precision))

			if _, err := batch.Del(generateGeoHashKey(key, hash, existPoint.Id)); err != nil {
				return err
			}
			if _, err := batch.Del(generateGeoHashBoxKey(key, marshalBox(box), existPoint.Id)); err != nil {
				return err
			}

			return nil
		}); err != nil {
		return false, err
	}

	return batch.Commit()
}

func GHAdd(key []byte, points []*pb.Point) (bool, error) {
	lock(LGeoHash, key)
	defer unlock(LGeoHash, key)

	metaKey := generateGeoHashMetaKey(key)
	precisionStr, err := store.Get(metaKey)
	if err != nil {
		return false, err
	}
	if precisionStr == nil || len(precisionStr) == 0 {
		return false, NotFoundGeoHashError
	}
	precision, err := strconv.ParseInt(string(precisionStr), 10, 32)
	if err != nil {
		return false, err
	}

	batch := store.NewBatch()
	defer batch.Close()

	for i := range points {
		point := points[i]
		value, err := proto.Marshal(point)
		if err != nil {
			return false, err
		}
		// find exist point
		existPointRaw, err := store.Get(generateGeoHashIdKey(key, point.Id))
		if err != nil {
			return false, err
		}
		if existPointRaw != nil && len(existPointRaw) > 0 {
			var existPoint pb.Point
			if err := proto.Unmarshal(existPointRaw, &existPoint); err != nil {
				return false, nil
			}
			hash, box := geohash.Encode(existPoint.Latitude, existPoint.Longitude, int(precision))

			if _, err := batch.Del(generateGeoHashIdKey(key, point.Id)); err != nil {
				return false, nil
			}
			if _, err := batch.Del(generateGeoHashKey(key, hash, point.Id)); err != nil {
				return false, err
			}
			if _, err := batch.Del(generateGeoHashBoxKey(key, marshalBox(box), point.Id)); err != nil {
				return false, err
			}
		}

		hash, box := geohash.Encode(point.Latitude, point.Longitude, int(precision))
		if _, err := batch.Set(generateGeoHashIdKey(key, point.Id), value); err != nil {
			return false, err
		}
		if _, err := batch.Set(generateGeoHashKey(key, hash, point.Id), value); err != nil {
			return false, err
		}
		if _, err := batch.Set(generateGeoHashBoxKey(key, marshalBox(box), point.Id), value); err != nil {
			return false, err
		}
	}

	return batch.Commit()
}

func GHRem(key []byte, points []*pb.Point) (bool, error) {
	lock(LGeoHash, key)
	defer unlock(LGeoHash, key)

	metaKey := generateGeoHashMetaKey(key)
	precisionStr, err := store.Get(metaKey)
	if err != nil {
		return false, err
	}
	if precisionStr == nil || len(precisionStr) == 0 {
		return false, NotFoundGeoHashError
	}
	precision, err := strconv.ParseInt(string(precisionStr), 10, 32)
	if err != nil {
		return false, err
	}

	batch := store.NewBatch()
	defer batch.Close()

	for i := range points {
		point := points[i]
		// find exist point
		existPointRaw, err := store.Get(generateGeoHashIdKey(key, point.Id))
		if err != nil {
			return false, err
		}
		if existPointRaw != nil && len(existPointRaw) > 0 {
			var existPoint pb.Point
			if err := proto.Unmarshal(existPointRaw, &existPoint); err != nil {
				return false, nil
			}
			hash, box := geohash.Encode(existPoint.Latitude, existPoint.Longitude, int(precision))

			if _, err := batch.Del(generateGeoHashIdKey(key, point.Id)); err != nil {
				return false, nil
			}
			if _, err := batch.Del(generateGeoHashKey(key, hash, point.Id)); err != nil {
				return false, err
			}
			if _, err := batch.Del(generateGeoHashBoxKey(key, marshalBox(box), point.Id)); err != nil {
				return false, err
			}
		}
	}

	return batch.Commit()
}

func GHGetBoxes(key []byte, latitude float64, longitude float64) ([]*pb.Point, error) {
	metaKey := generateGeoHashMetaKey(key)
	precisionStr, err := store.Get(metaKey)
	if err != nil {
		return nil, err
	}
	if precisionStr == nil || len(precisionStr) == 0 {
		return nil, NotFoundGeoHashError
	}
	precision, err := strconv.ParseInt(string(precisionStr), 10, 32)
	if err != nil {
		return nil, err
	}

	_, box := geohash.Encode(latitude, longitude, int(precision))

	res := make([]*pb.Point, 0)
	if err := store.Iterate0(generateGeoHashBoxPrefixKey(key, marshalBox(box)),
		func(key []byte, value []byte) error {
			var item pb.Point
			if err := proto.Unmarshal(value, &item); err != nil {
				return err
			}
			item.Distance = distance(&pb.Point{Latitude: latitude, Longitude: longitude}, &item)
			res = append(res, &item)
			return nil
		}); err != nil {
		return nil, err
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Distance < res[j].Distance
	})

	return res, nil
}

func GHGetNeighbors(key []byte, latitude float64, longitude float64) ([]*pb.Point, error) {
	metaKey := generateGeoHashMetaKey(key)
	precisionStr, err := store.Get(metaKey)
	if err != nil {
		return nil, err
	}
	if precisionStr == nil || len(precisionStr) == 0 {
		return nil, NotFoundGeoHashError
	}
	precision, err := strconv.ParseInt(string(precisionStr), 10, 32)
	if err != nil {
		return nil, err
	}

	neighbors := geohash.GetNeighbors(latitude, longitude, int(precision))

	res := make([]*pb.Point, 0)

	for i := range neighbors {
		if err := store.Iterate0(generateGeoHashPrefixKey(key, neighbors[i]),
			func(key []byte, value []byte) error {
				var item pb.Point
				if err := proto.Unmarshal(value, &item); err != nil {
					return err
				}
				item.Distance = distance(&pb.Point{Latitude: latitude, Longitude: longitude}, &item)
				res = append(res, &item)
				return nil
			}); err != nil {
			return nil, err
		}
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Distance < res[j].Distance
	})

	return res, nil
}

func GHCount(key []byte) (uint32, error) {
	count := uint32(0)
	if err := store.Iterate0(generateGeoHashIdPrefixKey(key),
		func(key []byte, value []byte) error {
			count++
			return nil
		}); err != nil {
		return 0, err
	}
	return count, nil
}

func GHMembers(key []byte) ([]*pb.Point, error) {
	res := make([]*pb.Point, 0)
	if err := store.Iterate0(generateGeoHashIdPrefixKey(key),
		func(key []byte, value []byte) error {
			var item pb.Point
			if err := proto.Unmarshal(value, &item); err != nil {
				return err
			}
			res = append(res, &item)
			return nil
		}); err != nil {
		return nil, err
	}
	return res, nil
}

func marshalBox(box *geohash.Box) string {
	return fmt.Sprintf("%32.32f:%32.32f:%32.32f:%32.32f", box.MinLat, box.MaxLat, box.MinLng, box.MaxLng)
}

func distance(one *pb.Point, two *pb.Point) (meter uint64) {
	earthRadius := 6378.1370
	d2r := math.Pi / 180
	dLong := (one.Longitude - two.Longitude) * d2r
	dLat := (one.Latitude - two.Latitude) * d2r
	a := math.Pow(math.Sin(dLat/2.0), 2.0) + math.Cos(one.Latitude*d2r)*math.Cos(two.Latitude*d2r)*math.Pow(math.Sin(dLong/2.0), 2.0)
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	d := earthRadius * c
	meter = uint64(d * 1000)
	return meter
}

func generateGeoHashMetaKey(key []byte) []byte {
	return []byte(fmt.Sprintf(geoHashMetaKeyTemplate, key))
}

// id key
func generateGeoHashIdPrefixKey(key []byte) []byte {
	return []byte(fmt.Sprintf(geoHashIdPrefixKeyTemplate, key))
}

func generateGeoHashIdKey(key []byte, id []byte) []byte {
	return []byte(fmt.Sprintf(geoHashIdKeyTemplate, key, id))
}

// hash key
func generateGeoHashPrefixKey(key []byte, hash string) []byte {
	return []byte(fmt.Sprintf(geoHashPrefixKeyTemplate, key, hash))
}

func generateGeoHashKey(key []byte, hash string, id []byte) []byte {
	return []byte(fmt.Sprintf(geoHashKeyTemplate, key, hash, id))
}

// box key
func generateGeoHashBoxPrefixKey(key []byte, box string) []byte {
	return []byte(fmt.Sprintf(geoHashBoxPrefixKeyTemplate, key, box))
}

func generateGeoHashBoxKey(key []byte, box string, id []byte) []byte {
	return []byte(fmt.Sprintf(geoHashBoxKeyTemplate, key, box, id))
}
