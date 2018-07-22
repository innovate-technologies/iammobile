$(document).ready(function () {
    $("#btnSubmit").on("click", function (e) {
        e.preventDefault();
        $("#btnSubmit").prop("disabled", true);
        $("#icon").attr("class", "fas fa-spinner fa-spin");
        $.get("/addme", gotResult);
    })  
});

function gotResult(data) {
    console.log(data)
    $("#btnSubmit").prop("disabled", false);
    $("#icon").attr("class", "fas fa-rocket");

    if (data.result == "ok") {
        $("#alertSuccess").removeClass("hidden");
        return
    }

    if (data.result == "failed") {
        $("#alertError").removeClass("hidden");
        return
    }

    if (data.result == "existing") {
        $("#alertExisting").removeClass("hidden");
        return
    }
}