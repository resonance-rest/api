Unofficial API & CDN for Wuthering Waves game on [**koyio.rest**](https://koyio.rest) - made with **Go**, stored data in **JSON** and deployed on [**Railway**](https://railway.app).

This is **not** a final version, I need people to complete the game data, if you are interested, write me on [**Telegram**](https://t.me/whosneksio).


# 🔗 API Reference

## Base URL

```http
  https://koyio.rest/
```

## Characters

#### Get character list

```http
  GET https://koyio.rest/characters
```

#### Get character's data

```http
  GET https://koyio.rest/characters/:name
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of a character   |

#### Get a character's image

```http
  GET https://koyio.rest/characters/:name/:type
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of a character   |
| `type`    | `string` | **Required** · `icon` or `portrait`  |

## Attributes

#### Get attribute list

```http
  GET https://koyio.rest/attributes
```

#### Get attribute's data

```http
  GET https://koyio.rest/attributes/:name
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of a attribute   |

#### Get a attribute's image

```http
  GET https://koyio.rest/attributes/:name/:type
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of a character   |
| `type`    | `string` | **Required** · `icon`                |
