const infoMessages = (inputName) => {
  if (inputName == "all") {
    $(".name-info-message").show();
    $(".mail-info-message").show();
    $(".title-info-message").show();
    $(".desc-info-message").show();
  } else if (inputName == "name") {
    $(".name-info-message").show();
  } else if (inputName == "email") {
    $(".mail-info-message").show();
  } else if (inputName == "title") {
    $(".title-info-message").show();
  } else if (inputName == "desc") {
    $(".desc-info-message").show();
  } else if (inputName == "success") {
    $(".name-info-message").hide();
    $(".mail-info-message").hide();
    $(".title-info-message").hide();
    $(".desc-info-message").hide();

    $("#success-modal").show();

    setTimeout(() => {
      $("#check").attr("checked", true);
    }, 1500);

    setTimeout(() => {
      $("#check").attr("checked", false);
      $("#success-modal").hide();
    }, 2200);
  }
};

const formSubmitEvent = (e) => {
  const nameControl =
    document.getElementById("name").value == "" ? true : false;
  const emailControl =
    document.getElementById("email").value == "" ? true : false;
  const titleControl =
    document.getElementById("title").value == "" ? true : false;
  const descControl =
    document.getElementById("desc").value == "" ? true : false;

  if (
    nameControl == true &&
    emailControl == true &&
    titleControl == true &&
    descControl == true
  ) {
    infoMessages("all");
  } else if (nameControl == true) {
    infoMessages("name");
  } else if (emailControl == true) {
    infoMessages("email");
  } else if (titleControl == true) {
    infoMessages("title");
  } else if (descControl == true) {
    infoMessages("desc");
  } else {
    var typeDefaultOption = document.getElementById("type").options[0].value;
    var priotiryDefaultOption = document.getElementById("priority").options[0]
      .value;

    $("#name").val("");
    $("#email").val("");
    $("#title").val("");
    $("#type").val(typeDefaultOption);
    $("#priority").val(priotiryDefaultOption);
    $("#desc").val("");

    infoMessages("success");
  }

  e.preventDefault();
};

const form = document.getElementById("issue-form");
form.addEventListener("submit", formSubmitEvent);
