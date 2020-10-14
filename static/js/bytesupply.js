$(document).ready(function() {
    setInterval(function() {
        var time = new Date();
        $("#doy").text(time.toDateString());
        $("#tod").text(time.toLocaleTimeString());
    }, 1000);

    $(".lang-es").hide();
});

function validateEntry(id, type) {
    Id = "#" + id;
    errId = Id + "Help";
    validEntry = true;
    $(errId).removeClass("d-block").addClass("d-none");

    if ($(Id).val().length == 0) {
        validEntry = false
    } else {
        switch(type) {
            case 'usphone':
                validEntry = (/^\D?(\d{3})\D?\D?(\d{3})\D?(\d{4})$/.test($(Id).val()));
                break;
            case 'email':
                validEntry = (/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+$/.test($(Id).val()));
                break;
            case 'text':
            default:
                validEntry = (/^[a-zA-Z0-9 \-\,\.&]+$/.test($(Id).val()));
        }
    }

    if (!validEntry) {
        $(errId).removeClass("d-none").addClass("d-block");
        $(Id).focus();
    }
}

function setLang(langId) {
    if (langId == "lang-es") {
        oldLangId = "lang-en"
    } else {
        oldLangId = "lang-es"
    }

    $("#"+langId).attr("src", "/static/img/" + oldLangId + ".png");
    $("#"+langId).attr("id", oldLangId);

    $("."+oldLangId).hide(function() {$("."+langId).show()});
}
