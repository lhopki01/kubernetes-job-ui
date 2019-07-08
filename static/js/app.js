// Enable any popovers
$(function () {
    $('[data-toggle="popover"]').popover()
    $('[data-toggle="tooltip"]').tooltip()
    var Moment = moment();
    var pld = $('#pageloaddate')
    pld.html(Moment.format('ddd D MMM HH:mm:ss'))
})