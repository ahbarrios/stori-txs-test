package main

func main() {
	// TODO The Source will implement some interface that always be delivered by the consumer (the func above) to
	// read from a file source (csv) and get transactions until io.EOF.
	// The implementation must be abstract enough to isolate the I/O ops from the Transaction processing flow

	// TODO Write a func that receive a Source: it will give Transactions as an output and visit that source for every accumulator.
	// It will have several Accumulators as input args that will be used as visitors for every processed transaction
	// extracted from the Source
	// The return type will be an error

	// TODO Write a func that receives an INPUT and return an HTML email template as result

	// TODO work on the SMTP server to send the email
}
