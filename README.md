# [PROJECT NAME]
[Project Description]

### Instrucciones para usar el Template (Borrar una vez completadas)

1. Clonar este repo y borrar el directorio .git
2. Substituir todas las apariciones de "template-go-api" con el nombre real del proyecto. Pueden usar alguna herramienta CLI o VSCode para hacerlo.
3. Iniciar un nuevo repo de GIT con: `git init`
4. Asignar el remoto origin correspondiente de Gitlab

Una vez tengan los elementos necesarios para deployar listos (stack deploy, cluster, vault, etc) es necesario modificar los siguientes archivos:

1. Modificar los archivos de vault en las rutas `vault/config.prod.hcl` y `vault/config.qa.hcl`
2. Modificar el pipeline de Gitlab. En especial el ultimo paso donde se actualiza el stack deploy.


## Requerimientos
- Go 1.16 (recomendado usar gvm o Docker)
- Make
- Docker

## Instrucciones para correr proyecto de forma local (host)

1. Primero es necesario preparar el archivo `.env` en la raiz del proyecto. El cual tiene el siguiente formato:

```env
ENVIRONMENT=dev
PORT=3000
JWT_SECRET=<YOUR_SECRET_HERE>
MONGO_URI=<YOUR_MONGO_URI>
MONGO_DATABASE=<YOUR_MONGO_DATABASE>
```
2. Luego haciendo uso de Make. Correr el siguiente comando:
```sh
make host-run
```
3. El proyecto deberia iniciar en el puerto 3000 o cualquier otro indicado por variables de entorno locales o en el archivo `.env`

## Instrucciones para correr proyecto usando Docker

1. Tambien es necesario contar con un archivo `.env` listo en la raiz del proyecto
2. Primero es necesario construir la imagen docker usando el siguiente comando make:
```sh
make docker-dev-build
```
3. Una vez construida la imagen, se puede ejecutar usando el siguiente comando:
```sh
make docker-dev-run
```
4.  El proyecto deberia iniciar en el puerto 3000 o cualquier otro indicado en el archivo `.env`

## Correr pruebas unitarias y de lintaje

Simplemente correr: 
```sh
make test
```
Para el linter:
```sh
make lint
```

## Debug server usando Docker

Es posible iniciar un debug server con docker. Primero hay que construir la imagen correspondiente:

```sh
make docker-debug-build
```
Luego iniciarla:
```sh
make docker-debug-run
```
Un servidor en modo debugger iniciara en el puerto 4000. Al cual pueden enlazar cualquier cliente de Debugger compatible con delv que posean. El de VSCode es bastante bueno.

## A la hora de deployar

Ya el proyecto cuenta con todos los pasos de Gitlab CICD configurados. En el momento de subir cualquier cambio al repo en una rama diferente a Master, se correran las verificaciones iniciales (Linter, Pruebas, Detector de secretos y vulnerabilidades). Una vez el pipeline finalice, el merge request (MR) puede ser mergeado con master. Una vez hecho esto, el resto de los pasos empezara a correr (compilacion, construccion de la imagen, etc)

No es posible subir cambios directamente a Master. Todo debe ser hecho a traves de un MR.