//index.js

var util = require("../../utils/util.js")

//获取应用实例
const app = getApp()

Page({
  data: {
    motto: 'Hello World',
    userInfo: {},
    hasUserInfo: false,
    canIUse: wx.canIUse('button.open-type.getUserInfo'),
    mimadaoh:'关于我们',
    // banner
    imgs_banner: [
      {
        url: '../../images/banner.jpg' ,
        name: '关于图链', 
        text: '图链记录你的点滴', 
      },
      { 
        url: '../../images/banner1.jpg',
        name: '关于我们', 
        text: '提供最优服务',  
       }
    ],

    // 链上图片列表 格式
    // {
    //   name: '链的Flag',
    //   timer: '2021-10-22 10:11:10', 
    //   tx: '0xcf2e558954f7962b329acafb4ae127c3972255e9afd4e2beaa156324b51269e0', 
    //   url:'http://61.160.212.59:8070/user/getPic?txHash=0xff821f56acadeabcd067333107d09c5a3c3d51f602f9be33184e7382f246bc2a',
    // }
    pictures: [],
    currentSwiper:0, 
    autoplay: true,
    github:'../../images/erwm.png',
    wixin:'../../images/erwm.png'
  },

  // banner自定义小圆点
  swiperChange: function (e) {
    this.setData({
      currentSwiper: e.detail.current
    })
  }, 

  //事件处理函数
  bindViewTap: function() {
    wx.navigateTo({
      url: '../logs/logs'
    })
  },

  //获取用户图片列表
  getPicList:  function(user, succ, fail, res) {

    //拼接网络请求字符串
    var content = "http://61.160.212.59:8070/user/getPicList?user="
    var url = content.concat(user)

    // 小程序顶部显示Loading
    wx.showNavigationBarLoading();
    wx.showLoading({
      title: "图片加载中...",
      icon: 'loading',
    })

    wx.request({
      url: url,
      header: {
        'content-type': 'application/json',
        'appVersion': '1.0.1'
      },
      method: 'get',
      success: res => {

        //打印返回信息
        // var resultStr = JSON.stringify(res.data)
        // console.log(resultStr)
        // var self = this
        // //刷新
        // self.setData({
        //   pictures: res.data.pictures
        // })
      },
      fail: err => {
        if (fail) fail(err);
      },
      complete: com => {
        wx.hideNavigationBarLoading();
        wx.hideLoading();

        //打印返回信息
        var resultStr = JSON.stringify(com.data)
        console.log(resultStr)

        var self = this
        //刷新
        self.setData({
          pictures: com.data.pictures
        })
      }
    })
  },

  //加载事件函数
  onLoad: function () {   
  
    this.getPicList('fxh7622')
    if (app.globalData.userInfo) {
      this.setData({
        userInfo: app.globalData.userInfo,
        hasUserInfo: true
      })
    } else if (this.data.canIUse){
      // 由于 getUserInfo 是网络请求，可能会在 Page.onLoad 之后才返回
      // 所以此处加入 callback 以防止这种情况
      app.userInfoReadyCallback = res => {
        this.setData({
          userInfo: res.userInfo,
          hasUserInfo: true
        })
      }
    } else {
      // 在没有 open-type=getUserInfo 版本的兼容处理
      wx.getUserInfo({
        success: res => {
          app.globalData.userInfo = res.userInfo
          this.setData({
            userInfo: res.userInfo,
            hasUserInfo: true
          })
        }
      })
    }
  },
  // getUserInfo: function(e) {
  //   console.log(e)
  //   app.globalData.userInfo = e.detail.userInfo
  //   this.setData({
  //     userInfo: e.detail.userInfo,
  //     hasUserInfo: true
  //   })
  // },

  // 获取滚动条当前位置
  onPageScroll: function (e) {
    console.log(e)
    if (e.scrollTop > 100) {
      this.setData({
        floorstatus: true
      });
    } else {
      this.setData({
        floorstatus: false
      });
    }
  },

  //回到顶部
  goTop: function (e) {  // 一键回到顶部
    if (wx.pageScrollTo) {
      wx.pageScrollTo({
        scrollTop: 0
      })
    } else {
      wx.showModal({
        title: '提示',
        content: '当前微信版本过低，无法使用该功能，请升级到最新微信版本后重试。'
      })
    }
  },
  // 跳转到列表页面
  jump_List:function(){
    wx.navigateTo({
      url: '../list/list',
    })
  }, 
})





