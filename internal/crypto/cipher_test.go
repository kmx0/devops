package crypto

import (
	"crypto/rsa"
	"errors"
	"math/big"
	"testing"

	"github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

func TestReadPublicKey(t *testing.T) {

	key := new(big.Int)
	key.SetString("105781886524689187163153549911195492915193724302230738399430902522659802525481893748497120346037910999652996206329794450089709625766571041615782758294694617911738725667324573331493085846935400096414579424605253154883072773420172514184034501115761383671033006915225979410464479856319422777863289477002628464979", 10)
	type wantStruct struct {
		pubkey *rsa.PublicKey
		err    error
	}

	tests := []struct {
		name     string
		filepath string
		want     wantStruct
	}{
		{
			name:     "not such file",
			filepath: "testpub1.pem",
			want: wantStruct{
				pubkey: nil,
				err:    errors.New("open testpub1.pem: no such file or directory"),
			},
		},
		{
			name:     "Wrong key",
			filepath: "testpub2.pem",
			want: wantStruct{
				pubkey: nil,
				err:    errors.New("asn1: structure error: tags don't match (16 vs {class:0 tag:0 length:159 isCompound:false}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} publicKeyInfo @3"),
			},
		},
		{
			name:     "Correct key",
			filepath: "testpub3.pem",
			want: wantStruct{
				pubkey: &rsa.PublicKey{
					N: key,
					E: 65537,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubkey, err := ReadPublicKey(tt.filepath)
			// PingDB(ctx context.Context, urlExample string) bool

			if pubkey != nil {
				assert.Equal(t, tt.want.pubkey.N.Cmp(pubkey.N), 0)
				assert.Equal(t, tt.want.pubkey.E, pubkey.E)
			} else {
				assert.Equal(t, tt.want.pubkey, pubkey)
			}
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			} else {
				assert.Equal(t, tt.want.err, err)
			}

		})
	}
}

func TestReadPrivateKey(t *testing.T) {

	key := new(big.Int)
	key.SetString("105781886524689187163153549911195492915193724302230738399430902522659802525481893748497120346037910999652996206329794450089709625766571041615782758294694617911738725667324573331493085846935400096414579424605253154883072773420172514184034501115761383671033006915225979410464479856319422777863289477002628464979", 10)
	type wantStruct struct {
		privkey *rsa.PrivateKey
		err     error
	}

	tests := []struct {
		name     string
		filepath string
		want     wantStruct
	}{
		{
			name:     "not such file",
			filepath: "testpriv1.pem",
			want: wantStruct{
				privkey: nil,
				err:     errors.New("open testpriv1.pem: no such file or directory"),
			},
		},
		{
			name:     "Wrong key",
			filepath: "testpriv2.pem",
			want: wantStruct{
				privkey: nil,
				err:     errors.New("asn1: structure error: tags don't match (16 vs {class:0 tag:0 length:604 isCompound:false}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} pkcs1PrivateKey @4"),
			},
		},
		{
			name:     "Correct key",
			filepath: "testpriv3.pem",
			want: wantStruct{
				privkey: &rsa.PrivateKey{
					PublicKey: rsa.PublicKey{
						N: key,
						E: 65537,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privkey, err := ReadPrivateKey(tt.filepath)
			// PingDB(ctx context.Context, urlExample string) bool

			if privkey != nil {
				assert.Equal(t, tt.want.privkey.PublicKey.N.Cmp(privkey.PublicKey.N), 0)
				assert.Equal(t, tt.want.privkey.PublicKey.E, privkey.PublicKey.E)
			} else {
				assert.Equal(t, tt.want.privkey, privkey)
			}
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			} else {
				assert.Equal(t, tt.want.err, err)
			}

		})
	}
}

func TestEncryptData(t *testing.T) {

	logrus.SetReportCaller(true)
	pubkey, _ := ReadPublicKey("testpub3.pem")
	pubkeyWrong, er := ReadPublicKey("testpub3.pem")
	logrus.Info(er)
	pubkeyWrong.N = big.NewInt(7878)
	type wantStruct struct {
		encryptredData []byte
		err            error
	}

	tests := []struct {
		name   string
		pubkey *rsa.PublicKey
		data   []byte
		want   wantStruct
	}{
		{
			name:   "test1",
			pubkey: pubkey,
			data:   []byte("test1"),
			want: wantStruct{
				encryptredData: []byte{68, 164, 63, 86, 28, 25, 55, 28, 32, 178, 34, 107, 104, 122, 50, 189, 168, 136, 19, 160, 209, 42, 209, 187, 49, 223, 137, 31, 191, 40, 41, 21, 222, 132, 155, 55, 102, 156, 196, 14, 48, 245, 52, 251, 184, 97, 175, 171, 137, 160, 104, 30, 250, 24, 27, 238, 117, 150, 106, 136, 76, 181, 32, 234, 100, 78, 248, 134, 36, 89, 0, 52, 128, 235, 47, 177, 74, 116, 80, 7, 213, 19, 217, 202, 152, 164, 239, 211, 34, 81, 170, 101, 131, 170, 46, 230, 177, 105, 41, 246, 124, 142, 14, 149, 96, 139, 251, 250, 35, 201, 184, 217, 197, 133, 22, 5, 159, 91, 95, 153, 190, 39, 1, 147, 247, 0, 181, 25},
				err:            nil,
			},
		},
		{
			name:   "test2",
			pubkey: pubkey,
			data:   []byte("test2"),
			want: wantStruct{
				encryptredData: []byte{39, 217, 230, 155, 3, 207, 48, 204, 215, 0, 240, 195, 75, 247, 227, 125, 128, 3, 91, 195, 137, 137, 224, 171, 58, 19, 13, 222, 231, 246, 158, 152, 82, 232, 71, 138, 28, 255, 199, 21, 213, 72, 125, 12, 117, 215, 189, 175, 125, 93, 212, 203, 90, 63, 246, 99, 26, 213, 183, 54, 143, 229, 74, 79, 228, 238, 58, 183, 28, 194, 173, 96, 150, 255, 185, 174, 47, 232, 163, 152, 78, 4, 2, 71, 82, 250, 225, 71, 93, 66, 66, 154, 183, 19, 189, 18, 145, 218, 16, 174, 213, 152, 151, 148, 222, 165, 10, 93, 197, 19, 37, 199, 216, 43, 13, 131, 223, 207, 63, 51, 255, 103, 58, 0, 88, 249, 244, 14},
				err:            nil,
			},
		},
		{
			name:   "test3 Wrong",
			pubkey: pubkeyWrong,
			data:   []byte("test"),
			want: wantStruct{
				err: errors.New("crypto/rsa: message too long for RSA public key size"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encdata, err := EncryptData(*tt.pubkey, tt.data)
			// PingDB(ctx context.Context, urlExample string) bool

			assert.Equal(t, len(tt.want.encryptredData), len(encdata))
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			} else {
				assert.Equal(t, tt.want.err, err)
			}

		})
	}
}

func TestDecryptData(t *testing.T) {

	logrus.SetReportCaller(true)
	privkey, _ := ReadPrivateKey("testpriv3.pem")
	privkeyWrong, er := ReadPrivateKey("testpriv3.pem")
	privkey.D = big.NewInt(12345)
	logrus.Info(er)
	// privkeyWrong.N = big.NewInt(7878)
	type wantStruct struct {
		data []byte
		err  error
	}

	tests := []struct {
		name           string
		privkey        *rsa.PrivateKey
		encryptredData []byte
		want           wantStruct
	}{
		{
			name:           "test1",
			privkey:        privkey,
			encryptredData: []byte{68, 164, 63, 86, 28, 25, 55, 28, 32, 178, 34, 107, 104, 122, 50, 189, 168, 136, 19, 160, 209, 42, 209, 187, 49, 223, 137, 31, 191, 40, 41, 21, 222, 132, 155, 55, 102, 156, 196, 14, 48, 245, 52, 251, 184, 97, 175, 171, 137, 160, 104, 30, 250, 24, 27, 238, 117, 150, 106, 136, 76, 181, 32, 234, 100, 78, 248, 134, 36, 89, 0, 52, 128, 235, 47, 177, 74, 116, 80, 7, 213, 19, 217, 202, 152, 164, 239, 211, 34, 81, 170, 101, 131, 170, 46, 230, 177, 105, 41, 246, 124, 142, 14, 149, 96, 139, 251, 250, 35, 201, 184, 217, 197, 133, 22, 5, 159, 91, 95, 153, 190, 39, 1, 147, 247, 0, 181, 25},
			want: wantStruct{
				data: []byte("test1"),
				err:  nil,
			},
		},
		{
			name:           "test2",
			privkey:        privkey,
			encryptredData: []byte{39, 217, 230, 155, 3, 207, 48, 204, 215, 0, 240, 195, 75, 247, 227, 125, 128, 3, 91, 195, 137, 137, 224, 171, 58, 19, 13, 222, 231, 246, 158, 152, 82, 232, 71, 138, 28, 255, 199, 21, 213, 72, 125, 12, 117, 215, 189, 175, 125, 93, 212, 203, 90, 63, 246, 99, 26, 213, 183, 54, 143, 229, 74, 79, 228, 238, 58, 183, 28, 194, 173, 96, 150, 255, 185, 174, 47, 232, 163, 152, 78, 4, 2, 71, 82, 250, 225, 71, 93, 66, 66, 154, 183, 19, 189, 18, 145, 218, 16, 174, 213, 152, 151, 148, 222, 165, 10, 93, 197, 19, 37, 199, 216, 43, 13, 131, 223, 207, 63, 51, 255, 103, 58, 0, 88, 249, 244, 14},
			want: wantStruct{
				data: []byte("test2"),
				err:  nil,
			},
		},
		{
			name:    "test3 Wrong",
			privkey: privkeyWrong,
			want: wantStruct{
				err: errors.New("crypto/rsa: decryption error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := DecryptData(tt.privkey, tt.encryptredData)
			// PingDB(ctx context.Context, urlExample string) bool

			assert.Equal(t, len(tt.want.data), len(data))
			if err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			} else {
				assert.Equal(t, tt.want.err, err)
			}

		})
	}
}
