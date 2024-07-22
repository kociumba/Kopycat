// let lineNumber = 0;

$(() => {
    const LogDisplay = (() => {
        let $codeElement;
        let lastLogContent = '';

        function init() {
            $('#log-section').load('/static/logDisplay.html', () => {
                $codeElement = $('#logs-hosted code');
                refreshLogs();
                setInterval(refreshLogs, 1000);
            });
        }

        function refreshLogs() {
            $.ajax({
                url: '/get-logs',
                method: 'GET',
                dataType: 'text',
                success: handleLogFetchSuccess,
                error: handleLogFetchError
            });
        }

        function handleLogFetchSuccess(data) {
            if ($codeElement && $codeElement.length) {
                if (data !== lastLogContent) {
                    const highlightedData = LogHighlighter.highlightLog(data || 'No logs available.');
                    // lineNumber ++;
                    $codeElement.html(highlightedData);
                    lastLogContent = data;
                }
            }
        }

        function handleLogFetchError(xhr, status, error) {
            console.error('Log fetch error:', error);
            if ($codeElement && $codeElement.length) {
                $codeElement.text('Failed to load logs. Please make sure Kopycat is running.');
            }
        }

        return { init };
    })();

    LogDisplay.init();
});

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

function submitFiles() {
    var originPath = $('#originPath').val();
    var destinationPath = $('#destinationPath').val();
    // var checkedDrive = $('input[name="drives-option"]:checked').val();

    console.log('File Path:', "originPath: ", originPath, ", destinationPath: ", destinationPath);
    // console.log('Checked Drive:', checkedDrive);

    // if (!checkedDrive) {
    //     console.error('No drive selected', checkedDrive);
    //     alert('Please select a drive.');
    //     return;
    // }

    if (!originPath) {
        console.error('No origin path entered');
        alert('Please enter an origin path.');
        return;
    }

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