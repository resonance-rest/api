# This is **not** a final version, I need people to complete the game data ⎯ if you are interested, write me on <a href="https://t.me/whosneksio" target="_blank">**Telegram**</a>.

Unofficial API & CDN for Wuthering Waves game on <a href="https://koyio.rest" target="_blank">**koyio.rest**</a> ⎯ made with **Go**, stored data in **JSON** and deployed on <a href="https://railway.app" target="_blank">**Railway**</a>.

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
