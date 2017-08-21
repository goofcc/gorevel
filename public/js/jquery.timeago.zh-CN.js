// Simplified Chinese
jQuery.timeago.settings.strings = {
  prefixAgo: null,
  prefixFromNow: null,
  suffixAgo: "前",
  suffixFromNow: "刚刚",
  seconds: "不到1分钟",
  minute: "1分钟",
  minutes: "%d分钟",
  hour: "1小时",
  hours: "%d小时",
  day: "1天",
  days: "%d天",
  month: "1个月",
  months: "%d月",
  year: "1年",
  years: "%d年",
  numbers: [],
  wordSeparator: "",
  formatter: function(prefix, words, suffix) { return [prefix, words, suffix].join(""); }
};