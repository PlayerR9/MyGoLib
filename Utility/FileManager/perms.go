package FileManager

import "io/fs"

const (
	// DP_OwnerOnly is the permission where only the owner can read, write, and execute.
	DP_OwnerOnly fs.FileMode = 0700

	// DP_OwnerRestrictOthers is the permission where the owner can read, write, and execute, and others can read and execute.
	DP_OwnerRestrictOthers fs.FileMode = 0755

	// DP_All is the permission where everyone can read, write, and execute.
	DP_All fs.FileMode = 0777

	// FP_OwnerOnly is the permission where only the owner can read and write.
	FP_OwnerOnly fs.FileMode = 0600

	// FP_OwnerRestrictOthers is the permission where the owner can read and write, and others can read.
	FP_OwnerRestrictOthers fs.FileMode = 0644

	// FP_All is the permission where everyone can read and write.
	FP_All fs.FileMode = 0666
)
