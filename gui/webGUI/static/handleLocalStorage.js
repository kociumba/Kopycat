$(() => {
    // Load previous state
    storageHandler.loadActiveMainTab();
    storageHandler.loadActiveFolderTab();
    storageHandler.loadUserInputs();
    storageHandler.loadSyncInterval();

    // Listen for state changes
    observers.observeActiveMainTab();
    observers.observeActiveFolderTab();
    observers.observeUserInputs();
    observers.observeSyncInterval();
});

const storageHandler = (() => {
    function saveActiveMainTab(activeMainTab) {
        localStorage.setItem('activeMainTab', activeMainTab);
    }

    function loadActiveMainTab() {
        const activeMainTab = localStorage.getItem('activeMainTab') || 'dashboard';
        $('.main-tabs button').removeClass('active-main-tab');
        $(`.main-tabs button[data-main-tab="${activeMainTab}"]`).addClass('active-main-tab');
        $('.main-tab-content').removeClass('active');
        $(`#${activeMainTab}`).addClass('active');
    }

    function saveActiveFolderTab(activeFolderTab) {
        localStorage.setItem('activeFolderTab', activeFolderTab);
    }

    function loadActiveFolderTab() {
        const activeFolderTab = localStorage.getItem('activeFolderTab') || 'tab1';
        $('.folder-tabs button').removeClass('active-folder-tab');
        $(`.folder-tabs button[data-folder-tab="${activeFolderTab}"]`).addClass('active-folder-tab');
        $('.folder-tab-content').removeClass('active');
        $(`#${activeFolderTab}`).addClass('active');
    }

    function saveUserInputs() {
        const userInputs = {};
        $('input[type="text"]').each(function() {
            userInputs[this.id] = $(this).val();
        });
        localStorage.setItem('userInputs', JSON.stringify(userInputs));
    }

    function loadUserInputs() {
        const userInputs = JSON.parse(localStorage.getItem('userInputs')) || {};
        for (const [id, value] of Object.entries(userInputs)) {
            $(`#${id}`).val(value);
        }
    }

    function saveSyncInterval(interval) {
        localStorage.setItem('syncInterval', interval);
    }

    function loadSyncInterval() {
        const savedInterval = localStorage.getItem('syncInterval');
        if (savedInterval) {
            $('#sync-interval').val(savedInterval);
        }
    }

    return {
        saveActiveMainTab,
        loadActiveMainTab,
        saveActiveFolderTab,
        loadActiveFolderTab,
        saveUserInputs,
        loadUserInputs,
        saveSyncInterval,
        loadSyncInterval
    };
})();

const observers = (() => {
    function observeActiveMainTab() {
        $('.main-tabs').on('click', 'button', (event) => {
            const newActiveMainTab = $(event.target).data('main-tab');
            storageHandler.saveActiveMainTab(newActiveMainTab);
            $('.main-tabs button').removeClass('active-main-tab');
            $(event.target).addClass('active-main-tab');
            $('.main-tab-content').removeClass('active');
            $(`#${newActiveMainTab}`).addClass('active');
        });
    }

    function observeActiveFolderTab() {
        $('.folder-tabs').on('click', 'button', (event) => {
            const newActiveFolderTab = $(event.target).data('folder-tab');
            storageHandler.saveActiveFolderTab(newActiveFolderTab);
            $('.folder-tabs button').removeClass('active-folder-tab');
            $(event.target).addClass('active-folder-tab');
            $('.folder-tab-content').removeClass('active');
            $(`#${newActiveFolderTab}`).addClass('active');
        });
    }

    function observeUserInputs() {
        $('input[type="text"]').on('input', () => {
            storageHandler.saveUserInputs();
        });
    }

    function observeSyncInterval() {
        $('#sync-interval').on('change', function() {
            const selectedInterval = $(this).val();
            storageHandler.saveSyncInterval(selectedInterval);
            // Assuming changeInterval() is defined elsewhere and handles the actual interval change
            // changeInterval();
        });
    }

    return {
        observeActiveMainTab,
        observeActiveFolderTab,
        observeUserInputs,
        observeSyncInterval
    };
})();
