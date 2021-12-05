var util = require("../../utils/util.js")

//获取应用实例
const app = getApp()

Page({

  data: {
    motto: 'Hello World',
    userInfo: {},
    hasUserInfo: false,

    // 链上图片列表 格式
    // {
    //   name: '链的Flag',
    //   timer: '2021-10-22 10:11:10', 
    //   tx: '0xcf2e558954f7962b329acafb4ae127c3972255e9afd4e2beaa156324b51269e0', 
    //   url:'http://61.160.212.59:8070/user/getPic?txHash=0xff821f56acadeabcd067333107d09c5a3c3d51f602f9be33184e7382f246bc2a',
    // }
    pictures: [],
  },

  //获取用户图片列表
  getPicList:  function(user, succ, fail, res) {

    //拼接网络请求字符串
    var content = "http://192.168.46.134:8070/user/getPicList?user="
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

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
  },
  //获取当前滑块的index
  bindchange: function (e) {
    const that = this;
    that.setData({
      currentData: e.detail.current
    })
  },
  //点击切换，滑块index赋值
  checkCurrent: function (e) {
    const that = this;

    if (that.data.currentData === e.target.dataset.current) {
      return false;
    } else {

      that.setData({
        currentData: e.target.dataset.current
      })
    }
  },
  // 外面的弹窗
  btn: function () {
    this.setData({
      showModal: true
    })
  },

  // // 禁止屏幕滚动
  // preventTouchMove: function() {
  // },

  // 弹出层里面的弹窗
  // ok: function () {
  //   this.setData({
  //     showModal: false
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
})
　　