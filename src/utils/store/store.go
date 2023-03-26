package store

import "main/src/utils"

// simple shared store

type Store struct {
	_contents  utils.SafeMap
	_menuState bool
}

func New() *Store {
	return &Store{
		_contents:  utils.NewSafeMap(),
		_menuState: false,
	}
}

// Add items into the resource store
//
// itemName should be formatted as described:
//	- (Format) --{package name acronym + comp/cmd}-{description of the item}
//	- (Example) --cccomp-channelIDs -> clearchannelcomp-channelIDs
func (s *Store) Bundle(item interface{}, itemName string) {
	s._contents.Write(itemName, item)
}

// Retrieve items from the resource store
//
// itemName is the name of the item given to bundle
func (s *Store) Acquire(itemName string) (data interface{}) {
	data, _ = s._contents.Read(itemName)
	return
}

// Remove items from the resource store
//
// itemName is the name of the item given to bundle
func (s *Store) Release(itemName string) {
	s._contents.Release(itemName)
}

// Modify the menu state
//
// State should be changed only if:
//	- ComponentType is SelectMenu
//	- Finished and cleaning up (set it to false)
func (s *Store) SetMenuState(state bool) {
	s._menuState = state
}

// See what the state of the menu is
func (s *Store) ViewMenuState() bool {
	return s._menuState
}
