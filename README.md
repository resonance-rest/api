<h1>I AM TAKING A BREAK, IF YOU WANT TO HELP, PULL REQUEST ON API</h1>

# 🔗 Unofficial API for Wuthering Waves on [**api.resonance.rest**](https://api.resonance.rest) - made with [**Go**](https://go.dev/), stored data in **JSON** and deployed on [**Railway**](https://railway.app).

### This is **not** a final version, I need people to complete the game data, if you are interested, write me on [**Telegram**](https://t.me/whosneksio).


# API Reference

## Base URL

```http
  https://api.resonance.rest/
```

## Characters

#### Get character list

```http
  GET https://api.resonance.rest/characters
```

#### Get character's data

```http
  GET https://api.resonance.rest/characters/:name
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of a character   |

#### Get a character's image

```http
  GET https://api.resonance.rest/characters/:name/:type
```

| Parameter | Type     | Description                                    |
| :-------- | :------- | :----------------------------------------------|
| `name`    | `string` | **Required** · name of a character             |
| `type`    | `string` | **Required** · `icon`, `portrait` or `circle`  |

## Emojis

#### Get character's emoji list

```http
  GET https://api.resonance.rest/characters/:name/emojis
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of a character   |

#### Get the emoji of a character

```http
  GET https://api.resonance.rest/characters/:name/emojis/:number
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of a character   |
| `number`  | `int`    | **Required** · number of a emoji     |

## Attributes

#### Get attribute list

```http
  GET https://api.resonance.rest/attributes
```

#### Get attribute's data

```http
  GET https://api.resonance.rest/attributes/:name
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of a attribute   |

#### Get a attribute's icon

```http
  GET https://api.resonance.rest/attributes/:name/icon
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of a attribute   |

## Weapons

#### Get weapon list

```http
  GET https://api.resonance.rest/weapons
```

#### Get weapons in type

```http
  GET https://api.resonance.rest/weapons/:type
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `type`    | `string` | **Required** · type of a weapon      |

#### Get a weapons's data

```http
  GET https://api.resonance.rest/weapons/:type:/:name
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `type`    | `string` | **Required** · type of a weapon      |
| `name`    | `string` | **Required** · name of a weapon      |

#### Get a weapons's image

```http
  GET https://api.resonance.rest/weapons/:type:/:name/icon
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `type`    | `string` | **Required** · type of a weapon      |
| `name`    | `string` | **Required** · name of a weapon      |

## Echoes

#### Get echoes list

```http
  GET https://api.resonance.rest/echoes
```

#### Get echo by name

```http
  GET https://api.resonance.rest/echoes/:name
```

| Parameter | Type     | Description                          |
| :-------- | :------- | :----------------------------------- |
| `name`    | `string` | **Required** · name of an echo       |

#### Get sonata effects list

```http
  GET https://api.resonance.rest/echoes/sonatas
```

#### Get sonata effect

```http
  GET https://api.resonance.rest/echoes/sonatas/:name
```

| Parameter | Type     | Description                             |
| :-------- | :------- | :-------------------------------------- |
| `name`    | `string` | **Required** · name of an sonata effect |

#### Get stats list

```http
  GET https://api.resonance.rest/echoes/stats
```

#### Get stat

```http
  GET https://api.resonance.rest/echoes/stats/:name
```

| Parameter | Type     | Description                             |
| :-------- | :------- | :-------------------------------------- |
| `name`    | `string` | **Required** · name of an substat       |

#### Get substats list

```http
  GET https://api.resonance.rest/echoes/substats
```

#### Get substat

```http
  GET https://api.resonance.rest/echoes/substats/:name
```

| Parameter | Type     | Description                             |
| :-------- | :------- | :-------------------------------------- |
| `name`    | `string` | **Required** · name of an substat       |


