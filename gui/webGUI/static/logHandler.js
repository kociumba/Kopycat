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
