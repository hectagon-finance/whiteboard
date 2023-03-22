Original block_data:
{
   "BlockHash":"dsahdasdhkadkjasdjsjkdhdkaskdjasjkd",
      "Instructions":[
         {
            "Id":"1",
            "C":"Create",
            "Data":{
               "Title":"Title1",
               "Desc":"Description1"
            }
         },
         {
            "Id":"2",
            "C":"Start",
            "Data":{
               "Id":"TEpWCQms",
               "EstDayToFinish": 2
            }
         },
         {
            "Id":"3",
            "C":"Stop",
            "Data":{
               "Id":"NOhpkzab",
               "Reason":"Reason3"
            }
         },
         {
            "Id":"4",
            "C":"Pause",
            "Data":{
               "Id":"ntZXIweo",
               "EstWaitDay":3
            }
         },
         {
            "Id":"5",
            "C":"Finish",
            "Data":{
               "Id":"TEpWCQms",
               "CongratMessage":"Description5"
         }
      }
   ]
}

=> encode Data field in Instructions to base64 to match with [ ]byte type:
{
  "BlockHash": "dsahdasdhkadkjasdjsjkdhdkaskdjasjkd",
  "Instructions": [
    {
      "Id": "1",
      "C": "Create",
      "Data": "eyJEZXNjIjoiRGVzY3JpcHRpb24xIiwiVGl0bGUiOiJUaXRsZTEifQ=="
    },
    {
      "Id": "2",
      "C": "Start",
      "Data": "eyJFc3REYXlUb0ZpbmlzaCI6MiwiSWQiOiJURXBXQ1FtcyJ9"
    },
    {
      "Id": "3",
      "C": "Stop",
      "Data": "eyJJZCI6Ik5PaHBremFiIiwiUmVhc29uIjoiUmVhc29uMyJ9"
    },
    {
      "Id": "4",
      "C": "Pause",
      "Data": "eyJFc3RXYWl0RGF5IjozLCJJZCI6Im50WlhJd2VvIn0="
    },
    {
      "Id": "5",
      "C": "Finish",
      "Data": "eyJDb25ncmF0TWVzc2FnZSI6IkRlc2NyaXB0aW9uNSIsIklkIjoiVEVwV0NRbXMifQ=="
    }
  ]
}

mem test original: 
[{"Id":"TEpWCQms","Title":"Title1","Desc":"Description1","Status":"JustCreated"},{"Id":"NOhpkzab","Title":"Title1","Desc":"Description1","Status":"JustCreated"},{"Id":"ntZXIweo","Title":"Title1","Desc":"Description1","Status":"JustCreated"},{"Id":"dhmSvwfI","Title":"Title1","Desc":"Description1","Status":"JustCreated"}]
