package main

func main() {
	// var token = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImEzck1VZ01Gdjl0UGNsTGE2eUYzekFrZnF1RSIsInR5cCI6IkpXVCIsIng1dCI6ImEzck1VZ01Gdjl0UGNsTGE2eUYzekFrZnF1RSJ9.eyJuYmYiOjE1NjQ1Njk3ODYsImV4cCI6MTU2NDU3MzM4NiwiaXNzIjoiaHR0cHM6Ly9zeXN0ZW10ZXN0LWN4aWQtZW1lYS1pZHAuYXp1cmV3ZWJzaXRlcy5uZXQiLCJhdWQiOlsiaHR0cHM6Ly9zeXN0ZW10ZXN0LWN4aWQtZW1lYS1pZHAuYXp1cmV3ZWJzaXRlcy5uZXQvcmVzb3VyY2VzIiwiY3hEb21haW5JbnRlcm5hbEFwaSIsInVzZXJNYW5hZ2VtZW50QXBpIl0sImNsaWVudF9pZCI6InV0Zm9yc2tlcmVuIiwic3ViIjoiZWVkODNjZGMtYTA4OS00MTMwLTk5NzUtODRjZDM1MTY5YTIxIiwiYXV0aF90aW1lIjoxNTY0NTY5Nzg2LCJpZHAiOiJGYWNlYm9vayIsImdpdmVuX25hbWUiOiJIw7JhIiwiZmFtaWx5X25hbWUiOlsiSHXhu7NuaCIsIkh14buzbmgiXSwiZW1haWwiOlsidGhhbmhob2EuYTFAZ21haWwuY29tIiwidGhhbmhob2EuYTFAZ21haWwuY29tIl0sInJvbGUiOiJVc2VyIiwic2NvcGUiOlsicHJvZmlsZSIsImN4cHJvZmlsZSIsIm9wZW5pZCIsImN4RG9tYWluSW50ZXJuYWxBcGkiLCJ1c2VyTWFuYWdlbWVudCJdLCJhbXIiOlsiZXh0ZXJuYWwiXX0.fn0KDbkAehNAsP94bBxzT2X3AEErCFUf3jSNHB2lyonMYHkhprWXoyoY-IyqThPvdEoyvT-yHJu3ZZEz6gfMNo6VhAzd4krksDcgCk3lhDpo8VzChl82txkCPw7vVHoJacQQgUCIl9SU3rvRP4a5XMCyi29Ss7fBWkuX33mgH8RXOdtXCxyojACEDcpDoO0shIj15qFziy83bcBEILJTXipdjvZkhvg1GWw5uYIKlvjaXv4IuWdsCp69HozRIjN-2O-1R1ALFZkKikyfOjUJi0nUhTU5FM27BISxc073zxHYjKqeaz9WmFAAjcYMBzzIo9vMTKrGU6vtElfGhgaSDQ"
	// parsedToken, err := jwt.ValidateToken(token)

	// if err == nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(parsedToken)

	// var token1 = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImEzck1VZ01Gdjl0UGNsTGE2eUYzekFrZnF1RSIsInR5cCI6IkpXVCIsIng1dCI6ImEzck1VZ01Gdjl0UGNsTGE2eUYzekFrZnF1RSJ9.eyJuYmYiOjE1NjQ1Njk3ODYsImV4cCI6MTU2NDU3MzM4NiwiaXNzIjoiaHR0cHM6Ly9zeXN0ZW10ZXN0LWN4aWQtZW1lYS1pZHAuYXp1cmV3ZWJzaXRlcy5uZXQiLCJhdWQiOlsiaHR0cHM6Ly9zeXN0ZW10ZXN0LWN4aWQtZW1lYS1pZHAuYXp1cmV3ZWJzaXRlcy5uZXQvcmVzb3VyY2VzIiwiY3hEb21haW5JbnRlcm5hbEFwaSIsInVzZXJNYW5hZ2VtZW50QXBpIl0sImNsaWVudF9pZCI6InV0Zm9yc2tlcmVuIiwic3ViIjoiZWVkODNjZGMtYTA4OS00MTMwLTk5NzUtODRjZDM1MTY5YTIxIiwiYXV0aF90aW1lIjoxNTY0NTY5Nzg2LCJpZHAiOiJGYWNlYm9vayIsImdpdmVuX25hbWUiOiJIw7JhIiwiZmFtaWx5X25hbWUiOlsiSHXhu7NuaCIsIkh14buzbmgiXSwiZW1haWwiOlsidGhhbmhob2EuYTFAZ21haWwuY29tIiwidGhhbmhob2EuYTFAZ21haWwuY29tIl0sInJvbGUiOiJVc2VyIiwic2NvcGUiOlsicHJvZmlsZSIsImN4cHJvZmlsZSIsIm9wZW5pZCIsImN4RG9tYWluSW50ZXJuYWxBcGkiLCJ1c2VyTWFuYWdlbWVudCJdLCJhbXIiOlsiZXh0ZXJuYWwiXX0.fn0KDbkAehNAsP94bBxzT2X3AEErCFUf3jSNHB2lyonMYHkhprWXoyoY-IyqThPvdEoyvT-yHJu3ZZEz6gfMNo6VhAzd4krksDcgCk3lhDpo8VzChl82txkCPw7vVHoJacQQgUCIl9SU3rvRP4a5XMCyi29Ss7fBWkuX33mgH8RXOdtXCxyojACEDcpDoO0shIj15qFziy83bcBEILJTXipdjvZkhvg1GWw5uYIKlvjaXv4IuWdsCp69HozRIjN-2O-1R1ALFZkKikyfOjUJi0nUhTU5FM27BISxc073zxHYjKqeaz9WmFAAjcYMBzzIo9vMTKrGU6vtElfGhgaSDQ"
	// parsedToken1, err := jwt.ValidateToken(token1)

	//fmt.Println(parsedToken1)
}
