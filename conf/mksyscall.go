/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019 WireGuard LLC. All Rights Reserved.
 */

package conf

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output zsyscall_windows.go dnsresolver_windows.go migration_windows.go storewatcher_windows.go
