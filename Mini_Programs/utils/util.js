//设置时间格式
const formatTime = date => {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()
  return [year, month, day].map(formatNumber).join('/') + ' ' + [hour, minute, second].map(formatNumber).join(':')
}

const formatNumber = n => {
  n = n.toString()
  return n[1] ? n : '0' + n
}

// 网络请求 url(), method(get、post)
const _netRequest = function(url, method, succ, fail, com) {

  // 小程序顶部显示Loading
  wx.showNavigationBarLoading();
  wx.showLoading({
    title: "图片加载中",
    icon: 'loading',
  })

  wx.request({
    url: url,
    header: {
      'content-type': 'application/x-www-form-urlencoded',
      'appVersion': appVersion||'' //小程序的版本号（可选,不填也不会报错）
    },
    method: method,
    success: res => {
      if (succ) succ(res);
    },
    fail: err => {
      if (fail) fail(err);
    },
    complete: com => {
      wx.hideNavigationBarLoading();
      wx.hideLoading();

      //设置数据
      var resultStr = JSON.stringify(com.data)
      console.log(resultStr)

      //设置
    }
  })
}

//获取用户图片列表
const getPicList = function(user, succ, fail, com) {

    //拼接网络请求字符串
    var content = "http://61.160.212.59:8070/user/getPicList?user="
    var content = http.concat(content, user)

    //网络请求
    _netRequest(url, "get", succ, fail, com)
}

//获取用户指定图片
const getPic = function(host, tx, succ, fail, com) {

  //拼接网络请求字符串
  var http = "http://"
  var content = "/user/getPic?txHash="
  var url = http.concat(host, content, tx)

  //网络请求
  _netRequest(url, "get", succ, fail)
}

/****测试****/
const getTestPic = function(succ, fail, com) {

  //拼接网络请求字符串
  var http = "http://61.160.212.59:8070"
  var content = "/user/getPic?txHash=0xcf2e558954f7962b329acafb4ae127c3972255e9afd4e2beaa156324b51269e0"
  var url = http.concat(http, content)

  //网络请求
  _netRequest(url, "get", succ, fail)
}

//透出函数
module.exports = {
  formatTime: formatTime,
  getPicList: getPicList,
  getPic:getPic,
  getTestPic:getTestPic
}
