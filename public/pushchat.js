$(function() {
    var chatboxIsScrolledDown = function() {
        var chatbox = $('#chatbox')
        return chatbox.prop('scrollHeight') - chatbox.scrollTop() == chatbox.outerHeight()
    }
    var chatBoxScrollDown = function () {
        var chatbox = $('chatbox')
        chatbox.scrollTop(chatbox.prop('scrollHeight'))
    }

})
