package secrets

import (
	"fmt"

	keychain "github.com/keybase/go-keychain"
)

var (
	ErrDuplicateItem = fmt.Errorf("Item already exists")
	ErrNotFound      = fmt.Errorf("secret not found in keyring")
)

type Keyring struct {
	kc keychain.Keychain
}

func NewKeyring(name string) (*Keyring, error) {
	return &Keyring{
		kc: keychain.NewWithPath(fmt.Sprintf("%s.keychain", name)),
	}, nil
}

func (k *Keyring) Set(service, username, pass string) error {
	_, _, err := k.Get(service)
	if err != nil && err != ErrNotFound {
		return err
	} else {
		return ErrDuplicateItem
	}

	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetAccount(username)
	item.SetData([]byte(pass))
	item.SetAccessible(keychain.AccessibleWhenUnlocked)
	item.UseKeychain(k.kc)
	err = keychain.AddItem(item)
	if err != nil {
		return err
	}
	if err == keychain.ErrorDuplicateItem {
		return ErrDuplicateItem
	}
	return nil
}

func (k *Keyring) Get(service string) (string, string, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnAttributes(true)
	query.SetReturnData(true)
	query.UseKeychain(k.kc)
	results, err := keychain.QueryItem(query)
	if err != nil {
		return "", "", err
	}
	if len(results) < 1 {
		return "", "", ErrNotFound
	}
	return results[0].Account, string(results[0].Data), nil
}

func (k *Keyring) Delete(service string) error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(service)
	item.UseKeychain(k.kc)
	return keychain.DeleteItem(item)
}
