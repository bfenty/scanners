# Scanners

This is an application to allow users to scan an ID and an order number into a SQL database using a barcode scanner. Manual entry is also viable, but it is primarily intended to be used on a Win CE scanner.

#Docker
I've designed this application primarily to run within Docker. To this end, I've included the dockerfile. use the code below in order to run in Docker:
>docker run -p 8080:8080 -e USER=username -e PASS=password -e SERVER=sqlserverip -e PORT=sqlserverport scanner
