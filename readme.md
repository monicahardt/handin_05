# Handin_05

To run the program you have to open up five different terminals

In the first three terminals, you have to start the three respective servers. They have to be started at port 5000, 5001 and 5002.

Feel free to copy the following three lines to each respective terminal to start the servers:

    go run server/server.go -port 5000
    go run server/server.go -port 5001
    go run server/server.go -port 5002

Now you have to start the frontend in the fourth terminal. The frontend has to run on it's own port. You can choose a port of your own liking, or just copy the following line to start the frontend:

    go run frontend.go -port 5010

Finally you have to start the client in the fifth terminal. The client has to run on it's own port and connect to the frontend port. Feel free to copy the following line to start the client:

    go run client/client.go -cPort 8181 -fPort 5010

In the client's terminal you can place a bid by writing 
    
    bid

to the terminal (followed by 'enter'). Then you can write the bid you want to place (on a new line!)

If you want to see the current result you can type 
    
    result

in the terminal

After 50 seconds the auction is closed and the client can only ask for the result.