package credmgrs

import (
	"fmt"
	"reflect"
	"sort"
)

const (
	userDataMapIamArnsKey = "iam_arns"
	userDataMapPwKey      = "password"
)

type DBKey struct {
	IAMArn string
	DBHost string
	DBPort string
	DBUser string
}

type DBCred struct {
	Host            string
	Port            string
	User            string
	Password        string
	AccessGrantedTo []string
}

// UserMap creates a map bearing the full key and value we'd see under
// a path comprised of the host and port.
func (c *DBCred) UserMap() map[string]interface{} {
	return map[string]interface{}{
		c.User: c.UserDataMap(),
	}
}

// UserDataMap is a convenience method for getting the map[string]interface{}
// that we'd write under a path comprised of the host and port, and under a key
// of the user's name.
func (c *DBCred) UserDataMap() map[string]interface{} {
	return map[string]interface{}{
		userDataMapPwKey:      c.Password,
		userDataMapIamArnsKey: c.AccessGrantedTo,
	}
}

func Sort(dbCreds []*DBCred) {
	sort.SliceStable(dbCreds, func(i, j int) bool {
		// Sort by DB host, port, then user.
		sortIBy := fmt.Sprintf("%s%s%s", dbCreds[i].Host, dbCreds[i].Port, dbCreds[i].User)
		sortJBy := fmt.Sprintf("%s%s%s", dbCreds[j].Host, dbCreds[j].Port, dbCreds[j].User)
		return sortIBy < sortJBy
	})
}

func toDbCreds(host, port string, secret map[string]interface{}) ([]*DBCred, error) {
	var dbCreds []*DBCred
	for dbUser, userDataMapIfc := range secret {
		userDataMap, ok := userDataMapIfc.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unable to convert userDataMapIfc to map, it's a %T", userDataMapIfc)
		}

		iamArns, err := getIamArns(userDataMap)
		if err != nil {
			return nil, err
		}

		var accessGrants []string
		for _, iamArn := range iamArns {
			accessGrants = append(accessGrants, iamArn)
		}

		password, err := getPassword(userDataMap)
		if err != nil {
			return nil, err
		}

		dbCreds = append(dbCreds, &DBCred{
			Host:            host,
			Port:            port,
			User:            dbUser,
			Password:        password,
			AccessGrantedTo: accessGrants,
		})
	}
	return dbCreds, nil
}

func getPasswordIfAuthorized(dbCreds []*DBCred, identity DBKey) (string, error) {
	for _, dbCred := range dbCreds {
		if dbCred.User != identity.DBUser {
			continue
		}
		authorized := false
		for _, iamArn := range dbCred.AccessGrantedTo {
			if iamArn == identity.IAMArn {
				authorized = true
				break
			}
		}
		if !authorized {
			return "", ErrNotAuthorized
		}
		return dbCred.Password, nil
	}
	return "", fmt.Errorf("username %s not found in stored credentials", identity.DBUser)
}

func deleteIfExists(toDel *DBCred, all []*DBCred) []*DBCred {
	for i, cred := range all {
		if !reflect.DeepEqual(cred, toDel) {
			continue
		}
		all = append(all[:i], all[i+1:]...)
		return all
	}
	return all
}

func getIamArns(userDataMap map[string]interface{}) ([]string, error) {
	iamArnsRaw, ok := userDataMap[userDataMapIamArnsKey]
	if !ok {
		return nil, fmt.Errorf("%s not found in %s", userDataMapIamArnsKey, userDataMap)
	}
	iamArns, ok := iamArnsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("could not convert %s to array, type is %T", iamArnsRaw, iamArnsRaw)
	}
	iamArnStrs := make([]string, len(iamArns))
	for i, iamArn := range iamArns {
		iamArnStr, ok := iamArn.(string)
		if !ok {
			return nil, fmt.Errorf("couldn't convert %s to string, it's a %T", iamArn, iamArn)
		}
		iamArnStrs[i] = iamArnStr
	}
	return iamArnStrs, nil
}

func getPassword(userDataMap map[string]interface{}) (string, error) {
	passwordRaw, ok := userDataMap[userDataMapPwKey]
	if !ok {
		return "", fmt.Errorf("password not found in %s", userDataMap)
	}
	password, ok := passwordRaw.(string)
	if !ok {
		return "", fmt.Errorf("could not convert %s to string, type is %T", passwordRaw, passwordRaw)
	}
	return password, nil
}

func toSecret(dbCreds []*DBCred) map[string]interface{} {
	secret := make(map[string]interface{})
	for _, dbCred := range dbCreds {
		secret[dbCred.User] = dbCred.UserDataMap()
	}
	return secret
}
