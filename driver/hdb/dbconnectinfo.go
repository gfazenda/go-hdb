// SPDX-FileCopyrightText: 2019-2021 Stefan Miller
//
// SPDX-License-Identifier: Apache-2.0

package hdb

type DBConnectInfo struct {
	DatabaseName string
	Host         string
	Port         int
	IsConnected  bool
}
