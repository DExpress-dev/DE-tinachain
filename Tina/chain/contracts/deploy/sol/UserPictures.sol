pragma solidity ^0.4.8;
pragma experimental ABIEncoderV2;

contract UserPictures {

    //全局变量，记录所有用户数量和图片数量;
    uint public userTotal;
    uint public picTotal;
    enum saveMode {IPFS, Tinachain}

    //保存图片的结构;
    struct pictures {
        bool hased;
        string picName;
        uint256 postion;
        saveMode mode;
        string uri;
        uint savetime;
    }

    //保存用户的结构;
    struct saveUsers{
        bool hased;
        string openId;  //1：微信小程序对应用户的openId。2：非微信小程序对应的用户名
        pictures[] picArray;
        mapping (bytes32 => pictures) picBook;
    }

    //保存用户的map;
    mapping (string => saveUsers) public usersBook;
    saveUsers[] public usersArray;

    constructor() public {
    }

    //添加图片
    function pictureAdd(string memory openId, string memory picName, uint256 postion, saveMode mode, string memory uri) public returns(bool){

        saveUsers storage user = usersBook[openId];
        if(user.hased == false){
            user.hased = true;
            user.openId = openId;
            userTotal++;
        }

        bytes32 orderID = keccak256(abi.encodePacked(picName, postion, mode, now));
        pictures storage savePic = user.picBook[orderID];
        if(savePic.hased == false){
            savePic.hased = true;
            savePic.picName = picName;
            savePic.postion = postion;
            savePic.mode = mode;
            savePic.uri = uri;
            savePic.savetime = now;
            user.picArray.push(savePic);
            user.picBook[orderID] = savePic;
            usersArray.push(user);
            usersBook[openId] = user;
            picTotal++;
            return (true);
        }
        return (false);
    }

    //得到指定用户所有图片信息
    function pictureGets(string memory openId) public view returns(string[] memory picNames, uint256[] memory postions, saveMode[] memory modes, 
        string[] memory uris, uint256[] memory savetimes) {

        saveUsers storage user = usersBook[openId];
        if(user.hased == true){
            uint256 picNum = user.picArray.length;
            picNames = new string[](picNum);
            postions = new uint256[](picNum);
            modes = new saveMode[](picNum);
            uris = new string[](picNum);
            savetimes = new uint256[](picNum);

            for(uint256 curIndex = 0; curIndex < picNum; curIndex++){
                pictures storage pic = user.picArray[curIndex];
                picNames[curIndex] = pic.picName;
                postions[curIndex] = pic.postion;
                modes[curIndex] = pic.mode;
                uris[curIndex] = pic.uri;
                savetimes[curIndex] = pic.savetime;
            }
            return(picNames, postions, modes, uris, savetimes);
        }
    }

    //得到指定用户所有图片信息
    // function getAllPics() public view returns(string[] memory openIds, string[] memory picNames, uint256[] memory postions, uint256[] memory modes, 
    //     string[] memory uris, uint256[] memory savetimes) {

    //     openIds = new string[](picTotal);
    //     picNames = new string[](picTotal);
    //     postions = new uint256[](picTotal);
    //     modes = new uint256[](picTotal);
    //     uris = new string[](picTotal);
    //     savetimes = new uint256[](picTotal);

    //     //轮询获取所有用户信息
    //     uint256 postion = 0;
    //     uint256 userIndex = 0;
    //     while(userIndex < usersArray.length){

    //         saveUsers storage user = usersArray[userIndex];
    //         for(uint256 picIndex = 0; picIndex < user.picArray.length; picIndex++){

    //             pictures storage picture = user.picArray[picIndex];
    //             openIds[postion] = user.openId;
    //             picNames[postion] = picture.picName;
    //             postions[postion] = picture.postion;
    //             modes[postion] = picture.mode;
    //             uris[postion] = picture.uri;
    //             savetimes[postion] = picture.savetime;
    //             postion++;
    //         }
    //         userIndex++;
    //     }
    //     return(openIds, picNames, postions, modes, uris, savetimes);
    // }

    //得到指定用户所有图片信息
    function getOpenIds() public view returns(string[] memory openIds) {

        openIds = new string[](usersArray.length);
        for(uint256 userIndex = 0; userIndex < usersArray.length; userIndex++){

            openIds[userIndex] = usersArray[userIndex].openId;
        }
        return(openIds);
    }
}