// SPDX-FileCopyrightText: 2014-2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/SAP/go-hdb/driver/common"
)

// DriverConn enhances a connection with go-hdb specific connection functions.
type DriverConn interface {
	ServerInfo() *common.ServerInfo
}
