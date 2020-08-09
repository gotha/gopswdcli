/**
 * The code below is based on zalando/go-keyring work
 * @author - https://github.com/zalando/go-keyring
 * @author - http://github.com/gotha/gopswdcli
 * @licence
 * The MIT License (MIT)
 *
 * Copyright (c) 2016 Zalando SE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package secrets

import (
	"fmt"

	"github.com/godbus/dbus"
)

var (
	// ErrNotFound is the expected error if the secret isn't found in the
	// keyring.
	ErrNotFound = fmt.Errorf("secret not found in keyring")
)

// KeyringLinux - dbus specific keyring for linux
type KeyringLinux struct {
	collection    string
	secretService *SecretService
}

func NewLinuxKeyring(collection string) (*KeyringLinux, error) {
	secretService, err := NewSecretService()
	if err != nil {
		return nil, fmt.Errorf("error creating secret service: %w", err)
	}
	return &KeyringLinux{
		collection:    collection,
		secretService: secretService,
	}, nil
}

func NewDefaultLinuxKeyring() (*KeyringLinux, error) {
	return NewLinuxKeyring("login")
}

func (k *KeyringLinux) Set(service, username, pass string) error {

	// open a session
	session, err := k.secretService.OpenSession()
	if err != nil {
		return err
	}
	defer k.secretService.Close(session)

	attributes := map[string]string{
		"service":  service,
		"username": username,
	}

	secret := NewSecret(session.Path(), pass)

	collection := k.secretService.GetSecretsCollection(k.collection)

	err = k.secretService.Unlock(collection.Path())
	if err != nil {
		return err
	}

	err = k.secretService.CreateItem(
		collection,
		fmt.Sprintf("Password for '%s' on '%s' ", service, username),
		attributes,
		secret,
	)
	if err != nil {
		return err
	}

	return nil
}

// findItem looksup an item by service and user.
func (k *KeyringLinux) findItem(key string) (dbus.ObjectPath, error) {

	collection := k.secretService.GetSecretsCollection(k.collection)
	search := map[string]string{
		"service": key,
	}

	err := k.secretService.Unlock(collection.Path())
	if err != nil {
		return "", err
	}

	results, err := k.secretService.SearchItems(collection, search)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", ErrNotFound
	}

	return results[0], nil
}

// Get gets a secret from the keyring given a key
func (k *KeyringLinux) Get(key string) (string, string, error) {

	item, err := k.findItem(key)
	if err != nil {
		return "", "", err
	}

	// open a session
	session, err := k.secretService.OpenSession()
	if err != nil {
		return "", "", err
	}
	defer k.secretService.Close(session)

	secret, err := k.secretService.GetSecret(item, session.Path())
	if err != nil {
		return "", "", err
	}

	attributes, err := k.secretService.GetSecretAttributes(item, session.Path())
	if err != nil {
		return "", "", err
	}

	username, exists := attributes["username"]
	if !exists {
		return "", "", fmt.Errorf("username attribute does not exist")
	}

	return username, string(secret.Value), nil
}

// Delete deletes a secret, identified by service & user, from the keyring.
func (k *KeyringLinux) Delete(key string) error {

	item, err := k.findItem(key)
	if err != nil {
		return err
	}

	return k.secretService.Delete(item)
}
