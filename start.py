import msgpack
import requests

'''data = {
    "Login": "tux",
    "Password": "12345678",
    "FirstName": "Tux",
    "LastName": "Torvalds",
    "Patronymic": "-",
}

query = msgpack.packb(data, use_bin_type=True)
response = requests.post('http://localhost:8080/api/user_create', data=query)

print(response.status_code)
print(response.content)'''
data = {
    "Token":       "CCmS1zwIRIukZwii31xOwrkrA2cz+fWA",
	"Id":        5,
    "Name":        "string",
	"Description": "string",
	"District":   "string",
	"Region":      "string",
	"Address":     "*string",
	"Type":        "*string",
	"State":       "*string",
	"Area":        100,
	"Owner":       "string",
	"Actual_user": "string",
	"Permissions": 1


}

query = msgpack.packb(data, use_bin_type=True)
response = requests.post('http://localhost:8080/api/object_change', data=query)

print(response.status_code)
print(msgpack.unpackb(response.content))
'''
data = {
    "Token": "CCmS1zwIRIukZwii31xOwrkrA2cz+fWA",
    "Name": "Main"
}

query = msgpack.packb(data, use_bin_type=True)
response = requests.post('http://localhost:8080/api/group_create', data=query)

print(response.status_code)
print(msgpack.unpackb(response.content))'''
