# ms_mod: Microsoft modules shipped with the Go toolset

This dir contains modules that are shipped with the Microsoft Go toolset.
It is essentially a vendor directory that takes priority over a project's local vendor directory.
Patches to the toolset enable this behavior if `GOEXPERIMENT` includes `xcryptobackendswap`.
