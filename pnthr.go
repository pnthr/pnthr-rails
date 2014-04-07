package main

import (
  "github.com/codegangsta/martini"
  "net/http"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "io/ioutil"
  "crypto/aes"
  "crypto/cipher"
  "encoding/base64"
  "fmt"
)

type Instance struct {
  Identifier string
  Secret string
  Password string
}

func main() {

  m := martini.Classic()

  /**
   * MongoDB Setup
   *
   * We are only connecting to a localhost for now
   * In the future we'll want to support multiple nodes which just means we need to pass
   * a comma separated list in the Dial() method
   */
  session, err := mgo.Dial("localhost")
  if err != nil {
    panic(err)
  }
  defer session.Close()

  /**
   * Optional. Switch the session to a monotonic behavior.
   *
   * Truthfully, i don't even know what the fuck this is
   */
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("pnthr_development").C("instances")

  /**
   * POST /
   *
   * All requests will come through the root, via post
   * Each request should have an app id and payload that has been encrypted with the app secret
   * We want to take this payload, decrypt it with the app secret
   * Once we have the raw payload, we encrypt first with the app password
   * Secondly with the encrypt with the app secret, for transport back to the requestor
   */
  m.Post("/", func(res http.ResponseWriter, req *http.Request) (int, string) {

    id := req.Header.Get("pnthr")

    /**
     * Failure: No App Key was passed
     *
     * If we don't have an api key in the header, then we can't fulfill this request
     */
    if len(id) == 0 {
      return 412, "Expected to find the 'pnthr' request header with your app id as the value, but none was found"
    }


    /**
     * Retrieve the application based on the passed App ID
     * If none is found we need to change the response status code
     */
    instance := Instance{}
    iv := []byte(id)[:aes.BlockSize]

    err = db.Find(bson.M{ "_id": bson.ObjectIdHex(id) }).One(&instance)

    if err != nil {
      panic(err)
    }

    /**
     * All has gone well with the request and app lookup
     * Let's decrypt the payload and then re-encrypt it
     */
    var (
      payload, err = ioutil.ReadAll(req.Body)
    )

    if err != nil {
      panic(err)
    }

    decoded := Base64Decode(string(payload))
    decrypted := make([]byte, len(string(decoded)))
    err = DecryptAES(decrypted, decoded, []byte(instance.Secret), iv)
    if err != nil {
      panic(err)
    }

    encrypted := make([]byte, len(string(decrypted)))
    err = EncryptAES(encrypted, decrypted, []byte(instance.Password), iv)
    if err != nil {
      panic(err)
    }

    fmt.Printf("Encrypted: %s\n", Base64Encode(encrypted))

    transport := make([]byte, len(string(encrypted)))
    err = EncryptAES(transport, encrypted, []byte(instance.Secret), iv)
    if err != nil {
      panic(err)
    }

    fmt.Printf("Transport: %s\n", Base64Encode(transport))

    decryptFirst := make([]byte, len(string(transport)))
    err = DecryptAES(decryptFirst, transport, []byte(instance.Secret), iv)
    if err != nil {
      panic(err)
    }
    fmt.Printf("Layer 1: %s\n", Base64Encode(decryptFirst))

    decryptSecond := make([]byte, len(string(decryptFirst)))
    err = DecryptAES(decryptSecond, decryptFirst, []byte(instance.Password), iv)
    if err != nil {
      panic(err)
    }

    fmt.Println("Layer 2: ", string(decryptSecond))


    return 200, Base64Encode(transport)
  })

  /**
   * POST /
   *
   * All requests will come through the root, via post
   * Each request should have an app id and payload that has been encrypted with the app secret
   * We want to take this payload, decrypt it with the app secret
   * Once we have the raw payload, we encrypt first with the app password
   * Secondly with the encrypt with the app secret, for transport back to the requestor
   */
  m.Post("/encrypt", func(res http.ResponseWriter, req *http.Request) (int, string) {

    id := req.Header.Get("pnthr")

    /**
     * Failure: No App Key was passed
     *
     * If we don't have an api key in the header, then we can't fulfill this request
     */
    if len(id) == 0 {
      return 412, "Expected to find the 'pnthr' request header with your app id as the value, but none was found"
    }

    /**
     * Retrieve the application based on the passed App ID
     * If none is found we need to change the response status code
     */
    instance := Instance{}
    iv := []byte(id)[:aes.BlockSize]

    err = db.Find(bson.M{ "_id": bson.ObjectIdHex(id) }).One(&instance)

    if err != nil {
      panic(err)
    }

    /**
     * All has gone well with the request and app lookup
     * Let's decrypt the payload and then re-encrypt it
     */
    var (
      payload, err = ioutil.ReadAll(req.Body)
    )

    if err != nil {
      panic(err)
    }

    encrypted := make([]byte, len(string(payload)))
    err = EncryptAES(encrypted, payload, []byte(instance.Secret), iv)
    if err != nil {
      panic(err)
    }

    fmt.Printf("Encrypting: %s\n", Base64Encode(encrypted))

    return 200, Base64Encode(encrypted)
  })


  m.Run()
}

func Base64Encode(b []byte) string {
    return base64.StdEncoding.EncodeToString(b)
}

func Base64Decode(s string) []byte {
    data, err := base64.StdEncoding.DecodeString(s)
    if err != nil {
        panic(err)
    }
    return data
}

func EncryptAES(dst, src, key, iv []byte) error {
  aesBlockEncryptor, err := aes.NewCipher([]byte(key))
  if err != nil {
    return err
  }
  aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncryptor, iv)
  aesEncrypter.XORKeyStream(dst, src)
  return nil
}

func DecryptAES(dst, src, key, iv []byte) error {
  aesBlockEncryptor, err := aes.NewCipher([]byte(key))
  if err != nil {
    return err
  }
  aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncryptor, iv)
  aesEncrypter.XORKeyStream(dst, src)
  return nil
}