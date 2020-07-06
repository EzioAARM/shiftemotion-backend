var jwt = require('jwt-simple');
var moment = require('moment');

var secret = 'mySecretTokenKey';

exports.handler = async (event, context, callback) => {
    var TOKEN = event.authorizationToken;
    try {
        var payload = jwt.decode(TOKEN, secret, false);
        if (payload) {
            if (payload.exp <= moment.unix()) {
                // vencio
                callback('Unauthorized');
            } else {
                // valido
                callback(null, generatePolicy('user', 'Allow', event.methodArn, payload.sub));
            }
        } else {
            callback('Unauthorized');
        }
    } catch (e) {
        callback('Unauthorized');
    }
    
}

// Help function to generate an IAM policy
var generatePolicy = function(principalId, effect, resource, user_id) {
    var authResponse = {}

    authResponse.principalId = principalId
    if (effect && resource) {
        var policyDocument = {}
        policyDocument.Version = '2012-10-17'
        policyDocument.Statement = []
        var statementOne = {}
        statementOne.Action = 'execute-api:Invoke'
        statementOne.Effect = effect
        statementOne.Resource = "*"
        policyDocument.Statement[0] = statementOne
        authResponse.policyDocument = policyDocument
    }

    // Optional output with custom properties of the String, Number or Boolean type.
    authResponse.context = {
        user_id_request: user_id,
    }
    return authResponse
}