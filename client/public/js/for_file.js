
jQuery(document).ready(function() {
    jQuery("#input-5").fileinput({showCaption: false});

    jQuery("#input-6").fileinput({
        showUpload: false,
        maxFileCount: 10,
        mainClass: "input-group-lg",
        showCaption: true
    });

    jQuery("#input-8").fileinput({
        mainClass: "input-group-md",
        showUpload: true,
        previewFileType: "image",
        browseClass: "btn btn-success",
        browseLabel: "Pick Image",
        browseIcon: "<i class=\"bi-image\"></i> ",
        removeClass: "btn btn-danger",
        removeLabel: "Delete",
        removeIcon: "<i class=\"bi-trash3\"></i> ",
        uploadClass: "btn btn-info",
        uploadLabel: "Upload",
        uploadIcon: "<i class=\"bi-upload\"></i> "
    });

    jQuery("#input-9").fileinput({
        previewFileType: "text",
        allowedFileExtensions: ["txt", "md", "ini", "text"],
        previewClass: "bg-warning",
        browseClass: "btn btn-primary",
        removeClass: "btn btn-secondary",
        uploadClass: "btn btn-secondary",
    });

    jQuery("#input-10").fileinput({
        showUpload: false,
        layoutTemplates: {
            main1: "{preview}\n" +
            "<div class=\'input-group {class}\'>\n" +
            "       {browse}\n" +
            "       {upload}\n" +
            "       {remove}\n" +
            "   {caption}\n" +
            "</div>"
        }
    });

    jQuery("#input-11").fileinput({
        maxFileCount: 10,
        allowedFileTypes: ["image", "video"]
    });

    jQuery("#input-12").fileinput({
        showPreview: false,
        allowedFileExtensions: ["zip", "rar", "gz", "tgz"],
        elErrorContainer: "#errorBlock"
    });
});

