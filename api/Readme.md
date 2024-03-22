# local run guide
-> just to to current dir
-> run go build ./...
-> go run cmd/main.go


# ec2 deployment guide
1. ssh to ec2 ssh -i "aws_local_connect.pem" ec2-user@ec2-18-215-149-126.compute-1.amazonaws.com
2. new setup |  install golang 1.22.0
```
curl 'https://dl.google.com/go/go1.22.0.linux-amd64.tar.gz' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:124.0) Gecko/20100101 Firefox/124.0' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Referer: https://go.dev/' -H 'Alt-Used: dl.google.com' -H 'Connection: keep-alive' -H 'Cookie: NID=512=PU9N6w_WO-IsmFDhPntIDgn-h4smMVcCgDfDvbx2E7M5U8XzdPnUmMwGoNZJv6DMU0sIJpd8_bgSs7wH1dqCk10z5PpNIKnTWNECcFkuXdwPj5o7HB67b8WS17_ZS3_-vtEucZv2ziXIt3guEGSi3aqGl-S_A-pBwxuV0wOVJU-skNrKv30WqOi9-0O91lu2gxB02VgEzbZJ78dBPHI-w-B9zwxYBmMV7NVXfeDNedtaWr-JwdrMzwg-NhnCw-AqfJs31Du4GfzLNNKu4ara4XYgdW1D3_bH7M4O5fO2KRDDAAN4Zuq9QWXHIHtub78oKveinJW-msfU1blwiS6WUpeaz6PLN_WH2whg4CbL5nm2lJX961O8rDxJE5_V9Km-ECC0MmPvxrElhAA8AQ2g4W2GEGRWM1fCaUfUpYAbzAQbj7zS47mF9Dg727EyNEuMf3O4KM71pukNuOYB62dBtVZSEGwrZmRzSTxE04AR827Ixjxezog2QpIrvepZC4RvKnbLGkUEprZXk_4DoyA356CXmBbfhxocny-6cuykr9ZQi2XIM1yndOwpl1PFs27yKdHx27WCurF_WUoz3IQp-xnK5xGJYP6AB9mAdPaAR-JEZVIy4MPVrbmEgxDTcgVtBZC_c1B3QjA_TkgiC7RStwmsgear2M4BQXIrrAw44Vau4iqUe__oCqJlCGkYgDayZGH7nMgKRSev8Pz4r-qDmz_nUECKFqk8UMlj1bEk_V5fqXxevCwQIU8uCAOl3Aq4_jnj5RCl8cRD-vk6G25y3TEHd53KjFuUg_BPG0KNw35x1YOKPxy_dp_rf52EdnDE844h7cuds1JgEzbjEl4Bu3auS8fsSXqjQve2GgvBYUnrFtow_eJPMdnu1sQUdULw7UCJKodaQCh7XB9CBCHUT8KbOLWkpnRAkZwA3Y0XsoXvLFpqU5ytl5aSs9NRO0w; 1P_JAR=2024-03-22-16; ANID=AHWqTUmsI7YblJKewBHsl0OllfjJmcRTtgmYgje3h3kE8CB21UrBiq4zeGe0ie5B; SID=g.a000hgh5ZstmhofJDiqaFJwf18V2aiaiT6zuP4YJ_kbxhfKrslGjDQFz71DNDLsAXDZX9ouIcQACgYKAZcSAQASFQHGX2MicLjVADBHwI8ciuT5rsmelBoVAUF8yKrK2t94wkaNHHlya3EeDVAV0076; __Secure-1PSID=g.a000hgh5ZstmhofJDiqaFJwf18V2aiaiT6zuP4YJ_kbxhfKrslGjnopwGAQQr1HI4K6H7L2XsAACgYKAdQSAQASFQHGX2MisiFrOoaiZIekSt3AvORkpRoVAUF8yKqHZexuuuTHu8y4674U3QCj0076; __Secure-3PSID=g.a000hgh5ZstmhofJDiqaFJwf18V2aiaiT6zuP4YJ_kbxhfKrslGjZX-mtEcL82-b9t6gNEH2UgACgYKAZwSAQASFQHGX2MilfNJT1y9vIfg1KzCnkWwBBoVAUF8yKrukKBso5W75ixY6M_AKsiI0076; HSID=AZBcSpY7vo8RRFwom; SSID=AnM8vkjh50kLyxI_e; APISID=37MbC5_Nm8PCrrzf/A4ELhWySo3_J_i6YE; SAPISID=8o4D5oLa29utdlhl/AAOKYpHzI98Vy_RJR; __Secure-1PAPISID=8o4D5oLa29utdlhl/AAOKYpHzI98Vy_RJR; __Secure-3PAPISID=8o4D5oLa29utdlhl/AAOKYpHzI98Vy_RJR; SIDCC=AKEyXzXNt4E4xcg1pqxfetiIGRkA0sDywQWNVJvmesuRrcMgbTPADNiBQnVvTCKMp6pf6i5xCneH; __Secure-1PSIDCC=AKEyXzUslS39jOY5VA1rqqcMjq2crqysfnADQNYnxgUgzYsyO_7bDpU5VR33Ap38OyZfM7qbx6uw; __Secure-3PSIDCC=AKEyXzVtJ47mOo3NfYOOAgTKVIxV1CUgoupCVGpwfsc0P-cHX_p0cMVkYo6v-IheMuSRlNmg8vU; __Secure-1PSIDTS=sidts-CjIB7F1E_DFhoM6kHSn2fhPnnyuRHIY-uvuFC58KXAjMV_fGdlIYLu578l7_IV4AV7uK5BAA; __Secure-3PSIDTS=sidts-CjIB7F1E_DFhoM6kHSn2fhPnnyuRHIY-uvuFC58KXAjMV_fGdlIYLu578l7_IV4AV7uK5BAA; AEC=Ae3NU9PAgVey9soiwbs3ddSyIoAh1tWV_YmT8O_qaZd2xkOGOcNhjAKz6Q' -H 'Upgrade-Insecure-Requests: 1' -H 'Sec-Fetch-Dest: document' -H 'Sec-Fetch-Mode: navigate' -H 'Sec-Fetch-Site: cross-site' -H 'Sec-Fetch-User: ?1' --output go1.22.0.linux-amd64.tar.gz
```
3. new setup | vim ~/.bash_profile
4. new setup | 
```
# GOROOT is the location where Go package is installed on your system
export GOROOT=/usr/local/go
# GOPATH is the location of your work directory
export GOPATH=$HOME/projects
# PATH in order to access go binary system wide
export PATH=$PATH:$GOROOT/bin
# Setup ENV for current sandbox
export ENV=dev
```
5. new setup | cd $HOME & mkdir projects
6. setup git ssh for ec2 and add in authorize in gihub repo in projects
7. clone git repo in projects
8. build, and run application
```
cd hastinapur/pkg go build ./...
cd hastinapur/api go build ./...
go build cmd/main.go
./main
```
9. run with ```nohup ./main &```
