# fbBot
facebook msg bot 

### 1. Just Deploy the same on Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Remember your heroku ID and app address. ex: `https://APP_ADDRESS.herokuapp.com/`

### 2. Paste Token to Heroku

Go to heroku dashboard, go to "Setting" -> "Config Variables".

- Add "Config Vars"
- Name -> "TOKEN"
- Value use  `PAGE_ACCESS_TOKEN` facebook app.


### 3 Persistent menu

```
curl -X POST -H "Content-Type: application/json" -d '{ 
    "get_started":{
        "payload":"GET_STARTED_PAYLOAD"
    }
}' "https://graph.facebook.com/v2.6/me/messenger_profile?access_token=PAGE_ACCESS_TOKEN"    
```

```
curl -X POST -H "Content-Type: application/json" -d '{
    "greeting":[
      {
        "locale":"default",
        "text":"Hello {{user_first_name}}!"
      }
    ]
}' "https://graph.facebook.com/v2.6/me/messenger_profile?access_token=PAGE_ACCESS_TOKEN"
```


```
curl -X POST -H "Content-Type: application/json" -d '{
    "persistent_menu":[
        {
            "locale":"default",
            "composer_input_disabled": false,
            "call_to_actions":[
                {
                    "title":"News",
                    "type":"nested",
                    "call_to_actions":[
                        {
                        "title":"Latest FoxNews",
                        "type":"postback",
                        "payload":"FOXNEWS"
                        }
                    ]
                },
                {
                    "title":"Post",
                    "type":"nested",
                    "call_to_actions":[
                        {
                        "title":"Latest Oziloo",
                        "type":"postback",
                        "payload":"OZILOO"
                        }
                    ]
                }
            ]
        }
    ]
}' "https://graph.facebook.com/v2.6/me/messenger_profile?access_token=PAGE_ACCESS_TOKEN"
```
