/**
 * @NApiVersion 2.x
 * @NScriptType Restlet
 */
define([ 'N/log', 'N/file', 'N/record', 'N/search' ], function (log, file, record, search) {
  var IS = search.Operator.IS;
  var ANYOF = search.Operator.ANYOF;
  var FOLDER = search.Type.FOLDER;
  var FILE = 'file'; // file does not exist in ns types

  var responses = {
    DELETED: {
      code: 200,
      status: 'DELETED'
    },
    SKIPPED: {
      code: 400,
      status: 'SKIPPED'
    },
    NOT_FOUND: {
      code: 404,
      status: 'NOT_FOUND'
    }
  };

  /**
   * Get the subfolders of a folder recursively
   * @param {string} id
   * @param {string} name
   * @param {string} [parent]
   * @returns {array<object>}
   */
  var getFolders = function getFolders(id, name, parent) {
    var path = parent ? parent + '/' + name : name;
    var folderIds = [ { id: id, path: path } ];
    search
      .create({
        type: FOLDER,
        filters: [ [ 'parent', IS, id ] ],
        columns: [ 'name' ]
      })
      .run()
      .each(function (res) {
        folderIds = [].concat(folderIds, getFolders(res.id, res.getValue('name'), path));
      });

    return folderIds;
  };

  /**
   * Find the id of an item by path
   * @param {string} path
   * @returns {array<object>}
   */
  var findIdsByName = function findIdsByName(path) {
    try {
      // check if file
      var item = file.load({ id: path });
      if (item) {
        log.debug('file found', JSON.stringify({ id: item.id }));
        return [
          {
            status: responses.DELETED.status,
            code: responses.DELETED.code,
            id: item.id,
            message: 'Successfully deleted "' + path + '"',
            path: path,
            type: FILE
          }
        ];
      }
    } catch (_) {
      // file not found
    }

    var name = path.replace(/^.*[\\\/]/, '');
    var folder = path.replace('/' + name, '');
    var parent = folder.split('/').pop();

    // check if folder
    var result = search
      .create({
        type: FOLDER,
        filters: [ [ 'name', IS, name ], 'and', [ 'parent', ANYOF, parent ] ],
        columns: [ 'name', 'parent' ]
      })
      .run()
      .getRange({ start: 0, end: 2 })
      .filter(function (res) {
        return res.getText('parent') === parent;
      });

    log.debug('search results', JSON.stringify(result));

    if (!result || !result.length) {
      return [
        {
          status: responses.NOT_FOUND.status,
          code: responses.NOT_FOUND.code,
          message: '"' + path + '" not found',
          path: path
        }
      ];
    }
    if (result.length > 1) {
      return [
        {
          status: responses.SKIPPED.status,
          code: responses.SKIPPED.code,
          message: 'Multiple results for "' + path + '"',
          path: path,
          result: result,
          type: FOLDER
        }
      ];
    }
    if (result.length === 1) {
      var id = result[0].id;
      var folders = getFolders(id, path).reverse();
      var files = [];
      search
        .create({
          type: FILE,
          filters: [ [ FOLDER, IS, id ] ],
          columns: [ 'name' ]
        })
        .run()
        .each(function (res) {
          return files.push({ id: res.id });
        });

      log.debug('folders', JSON.stringify(folders));
      log.debug('files', JSON.stringify(files));

      return [].concat(
        files.map(function (res) {
          var id = res.id;
          var path = file.load({ id: id }).path;
          return {
            status: responses.DELETED.status,
            code: responses.DELETED.code,
            message: 'Successfully deleted ' + FILE + ' "' + path + '"',
            type: FILE,
            id: id,
            path: path
          };
        }),
        folders.map(function (res) {
          var id = res.id;
          var path = res.path;
          return {
            status: responses.DELETED.status,
            code: responses.DELETED.code,
            message: 'Successfully deleted ' + FOLDER + ' "' + path + '"',
            type: FOLDER,
            id: id,
            path: path
          };
        }));
    }
  };

  /**
   * Delete the items specified as path strings in an array
   * @param {array|object} items
   */
  var deleteItems = function deleteItems(items) {
    var itemArray = items instanceof Array ? items : [ items ];
    if (!items.length) throw new Error('No items to delete');

    var reduced = itemArray.reduce(function (red, remoteFilePath) {
      var res = findIdsByName(remoteFilePath);
      return {
        del: [].concat(red.del, res.filter(function (res) {
          return res.id;
        })),
        err: [].concat(red.err, res.filter(function (res) {
          return !res.id;
        }))
      };
    }, { del: [], err: [] });

    var del = reduced.del;
    var err = reduced.err;

    var res = del.reduce(function (red, item) {
      try {
        switch (item.type) {
            case FILE:
              file.delete({ id: item.id });
              break;
            case FOLDER:
            default:
              record.delete({ type: item.type, id: item.id });
              break;
        }
        return {
          unsuccessful: red.unsuccessful,
          successful: [].concat(red.successful, [ item ])
        };
      } catch (error) {
        log.debug({ title: '"delete" error', details: JSON.stringify(error) });
        return {
          successful: red.successful,
          unsuccessful: [].concat(red.unsuccessful, [
            {
              id: item.id,
              type: item.type,
              path: item.path,
              code: 500,
              status: 'DELETE_ERROR',
              message: error.message
            }
          ])
        };
      }
    }, {
      successful: [],
      unsuccessful: err
    });

    log.debug({ title: '"delete" result', details: JSON.stringify(res) });
    return res;
  };

  return {
    // ROUTER
    get: function get() {
      return {
        code: 200,
        status: 'RESTLET_EXISTS',
        message: new Date().getTime().toString().substring(0, 10) + ' | The suitesync restlet is setup and healthy',
      };
    },
    post: function post(res) {
      var action = res.action;
      try {
        switch (action) {
            case 'delete':
              return deleteItems(res.items, log, record);
            default:
              throw new Error('Invalid action "' + action + '"');
        }
      } catch (err) {
        log.error({
          title: err.message,
          details: JSON.stringify(err)
        });
        return {
          error: {
            code: 500,
            status: 'MISC_ERROR',
            message: 'Error processing "' + action + '", ' + err.message,
            datain: res
          }
        };
      }
    }
  };
});
