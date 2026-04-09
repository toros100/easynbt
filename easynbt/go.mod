module github.com/toros100/easynbt/easynbt

go 1.26.1

require (
	github.com/google/go-cmp v0.7.0
	github.com/toros100/easynbt/nbt v0.1.0
	golang.org/x/tools v0.43.0
)

replace github.com/toros100/easynbt/nbt => ../nbt/

replace github.com/toros100/easynbt/nbt/nbtcmp => ../nbt/nbtcmp/

require (
	golang.org/x/mod v0.34.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
)
