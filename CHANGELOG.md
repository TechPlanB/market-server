
<a name="v1.0.10"></a>
## [v1.0.10](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.9...v1.0.10) (2021-10-12)


<a name="v1.0.9"></a>
## [v1.0.9](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.8...v1.0.9) (2021-10-12)

### Feat

* release v1.0.9


<a name="v1.0.8"></a>
## [v1.0.8](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.7...v1.0.8) (2021-10-12)

### Feat

* release v1.0.8

### Fix

* remove properties for nft token info


<a name="v1.0.7"></a>
## [v1.0.7](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.6...v1.0.7) (2021-10-11)

### Feat

* release v1.0.7
* add auction id field in artwork list
* add artwork list nft_address field
* add nft sales and artworks info properties field
* add multiple table query
* sales info logic
* add drops info category name and collection name logic
* add status code in sales
* artwork list logic and model
* add sales_list model and logic
* list logic
* basic framework of list logic
* basic framework of get info logic
* basic framework of get info logic
* generate sql model and types
* generate handlers
* redefine api
* sql model generated

### Fix

* price to fixed price
* rename nft address field
* artwork info query
* fix drops info query
* fix sales and artworks info query
* revise handlers
* add cache
* rename handler

### Merge Requests

* Merge branch 'dev/zhangyu' into 'develop'
* Merge branch 'dev/zhangyu' into 'develop'
* Merge branch 'dev/zhangyu' into 'develop'
* Merge branch 'feature/sql_define' into 'develop'
* Merge branch 'feature/sql_define' into 'develop'


<a name="v1.0.6"></a>
## [v1.0.6](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.5...v1.0.6) (2021-09-20)

### Feat

* release v1.0.6
* add more field for nft info api


<a name="v1.0.5"></a>
## [v1.0.5](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.4...v1.0.5) (2021-09-13)


<a name="v1.0.4"></a>
## [v1.0.4](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.3...v1.0.4) (2021-09-13)

### Feat

* release v1.0.4


<a name="v1.0.3"></a>
## [v1.0.3](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.2...v1.0.3) (2021-09-10)

### Feat

* release v1.0.3

### Fix

* fix deploy prod

### Improvement

* change develop trigger for ci

### Refactor

* change market go


<a name="v1.0.2"></a>
## [v1.0.2](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.1...v1.0.2) (2021-09-09)

### Feat

* release v1.0.2

### Fix

* fix gitlab ci
* fix gitlab ci yaml

### Improvement

* change image url for nft model

### Refactor

* change market.go


<a name="v1.0.1"></a>
## [v1.0.1](https://gitlab.ziggurat.cn:10000/nftm/market-server/compare/v1.0.0...v1.0.1) (2021-09-09)

### Docs

* change changelog

### Fix

* fix gitlab ci yaml


<a name="v1.0.0"></a>
## v1.0.0 (2021-09-09)

### Build

* change image build download envsubst

### Ci

* add prod dpeloyment
* remove clean all image for ci
* add clean image script
* add gitlab ci deployment script
* add gitlab ci yaml file

### Feat

* release v1.0.0
* add transactions api support
* add id_in_contract for nft info api
* add artwork detail info api
* add name,description for nft token
* add artworks api support
* add contract address for nft model
* add cover url for nft list api
* add nft count support for nft list api
* add nfts list api support
* add entrypoint.sh for docker startup
* add docker file support
* add get nft info api support

### Fix

* fix artwork info api name bug
* fix pagination for nft list
* fix nft info response define
* fix market api define
* fix docker command
* fix envsubst bin copy failed
* fix docker build error
* fix image build bug
* fix image build error

### Improvement

* add category filter for artworks
* add new nft sql script
* add categories api support
* add is_set for nft list api
* add description response for nft info
* add minio output url replace
* add nft files for nft info api
* refine download envsubst binary

### Refactor

* change docker build method
* change market.go

