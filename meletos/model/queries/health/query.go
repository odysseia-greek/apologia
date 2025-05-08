package health

import "fmt"

func Health() string {
	return fmt.Sprintf(`{
	health {
			time
			healthy
			version
			services{
				name
				version
				healthy
				databaseInfo{
					healthy
					serverName
					serverVersion
					clusterName
				}
			}
		}
	}
`)
}
