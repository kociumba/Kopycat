$(document).ready(function(){
    // Main tab switching
    $('.main-tabs button').click(function(){
        const tabId = $(this).data('main-tab');
        $('.main-tab-content').removeClass('active');
        $('#' + tabId).addClass('active');
        $('.main-tabs button').removeClass('active-main-tab');
        $(this).addClass('active-main-tab');
    });

    // Folder tab switching
    $('.folder-tabs button').click(function(){
        const tabId = $(this).data('folder-tab');
        $('.folder-tab-content').removeClass('active');
        $('#' + tabId).addClass('active');
        $('.folder-tabs button').removeClass('active-folder-tab');
        $(this).addClass('active-folder-tab');
    });
});

let selectedFolderPath = '';

// The new submiter needs to be this couse of the tabs 
function submitFiles() {
    // Find the active tab
    var activeTab = $('.folder-tab-content.active');
    
    // Get the input values from the active tab
    var originPath = activeTab.find('input[id^="originPath"]').val();
    var destinationPath = activeTab.find('input[id^="destinationPath"]').val();
    
    console.log('File Path:', "originPath: ", originPath, ", destinationPath: ", destinationPath);

    if (!originPath) {
        console.error('No origin path entered');
        alert('Please enter an origin path.');
        return;
    }

    // Check if we're in the "Mirror on Drive" tab
    if (activeTab.attr('id') === 'tab2') {
        var checkedDrive = activeTab.find('input[name="drives-option"]:checked').val();
        console.log('Checked Drive:', checkedDrive);

        if (!checkedDrive) {
            console.error('No drive selected');
            alert('Please select a drive.');
            return;
        }

        // Use the checked drive as the destination path
        destinationPath = checkedDrive;
    }

    // Call the addfolder function with the gathered information
    addfolder(originPath, destinationPath);
}

/** 
 * @param {string} originPath - the origin path of the folder to be added
 * @param {string} destinationPath - the destination path of the folder to be added
 */
function addfolder(originPath, destinationPath) {
    $.ajax({
        type: 'POST',
        url: '/add-folder',
        contentType: 'application/json',
        data: JSON.stringify({ origin: originPath, destination: destinationPath }),
        success: function (response) {
            console.log('Full path from backend:', response.fullPath);
            alert('Folder added to sync: ' + response.fullPath);
        },
        error: function (jqXHR, textStatus, errorThrown) {
            console.error('Error resolving folder path:', textStatus, errorThrown);
            console.log('Response:', jqXHR.responseText);
            alert('Error resolving folder path: ' + textStatus);
        }
    });
}

/**
 * @param {string} originPath - the origin path of the folder to be deleted
 * @param {string} destinationPath - the destination path of the folder to be deleted
 */
function deleteTarget(originPath, destinationPath) {
    // Implement your delete logic here
    console.log("Delete target:", originPath, "->", destinationPath);

    $.ajax({
        type: 'POST',
        url: '/delete-folder',
        contentType: 'application/json',
        data: JSON.stringify({ origin: originPath, destination: destinationPath }),
        success: function (response) {
            console.log('Full path from backend:', response.fullPath);
            alert('Folder deleted from sync: ' + response.fullPath);
        },
        error: function (jqXHR, textStatus, errorThrown) {
            console.error('Error resolving folder path:', textStatus, errorThrown);
            console.log('Response:', jqXHR.responseText);
            alert('Error resolving folder path: ' + textStatus);
        }
    })
}

function changeInterval() {
    var selectedOption = parseInt($("#sync-interval").val());
    $.ajax({
        type: 'POST',
        url: '/set-new-interval',
        contentType: 'application/json',
        data: JSON.stringify({ interval: selectedOption }),
        success: function (response) {
            console.log('Interval set to:', selectedOption);
        },
        error: function (jqXHR, textStatus, errorThrown) {
            console.error('Error setting interval:', textStatus, errorThrown);
        }
    });
}