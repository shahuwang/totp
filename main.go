package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

func main() {
	var key string
	flag.StringVar(&key, "key", "", "The key to generate two step code")
	var step, length int64
	flag.Int64Var(&step, "step", 30, "The step length")
	flag.Int64Var(&length, "len", 6, "The length of the code")
	flag.Parse()
	if key == "" {
		// read the default key in ~/.totp/key
		user, err := user.Current()
		if err != nil {
			fmt.Println(err)
			return
		}
		fp := fmt.Sprintf("%s/.totp/key", user.HomeDir)
		content, err := ioutil.ReadFile(fp)
		if err != nil || len(content) == 0 {
			fmt.Println(err)
			fmt.Println("You should put your two step key in this file")
			return
		}
		key = string(content)
		key = strings.TrimSpace(key)
	}
	clip(key, length, step)
}

func clip(key string, length, step int64) {
	code, left, err := totp(key, length, step)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The code is %s, %d seconds left, the code is copied to clipboard\n", code, left)
	path, err := exec.LookPath("xsel")
	if err != nil {
		fmt.Println(err)
		fmt.Println("May be you should install xsel by apt-get install xsel")
		return
	}
	cmd := exec.Command(path, "-b", "-i")
	cmd.Stdin = strings.NewReader(code)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func totp(key string, length, step int64) (string, int64, error) {
	k, _ := base32.StdEncoding.DecodeString(key)
	timebuf := new(bytes.Buffer)
	timestamp := time.Now().Unix()
	err := binary.Write(timebuf, binary.BigEndian, timestamp/step)
	left := step - timestamp%step
	if err != nil {
		fmt.Println(err)
		return "", 0, err
	}
	hash := hmac.New(sha1.New, []byte(k))
	hash.Write(timebuf.Bytes())
	v := hash.Sum(nil)
	o := v[len(v)-1] & 0xf
	c := (int32(v[o]&0x7f)<<24 | int32(v[o+1])<<16 | int32(v[o+2])<<8 | int32(v[o+3])) % 1000000000
	return fmt.Sprintf("%010d", c)[10-length : 10], left, nil
}
