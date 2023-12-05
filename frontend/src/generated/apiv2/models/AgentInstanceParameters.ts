/* tslint:disable */
/* eslint-disable */
/**
 * Cryptocat API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.3
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { AgentInstanceState } from './AgentInstanceState';
import {
    AgentInstanceStateFromJSON,
    AgentInstanceStateFromJSONTyped,
    AgentInstanceStateToJSON,
} from './AgentInstanceState';

/**
 * 
 * @export
 * @interface AgentInstanceParameters
 */
export interface AgentInstanceParameters {
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceParameters
     */
    endpoint: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceParameters
     */
    buildId: string;
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceParameters
     */
    configVersion: string;
    /**
     * 
     * @type {AgentInstanceState}
     * @memberof AgentInstanceParameters
     */
    state: AgentInstanceState;
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceParameters
     */
    tlsCertificateId: string;
    /**
     * 
     * @type {boolean}
     * @memberof AgentInstanceParameters
     */
    tlsCertificateSignedByPublicCa: boolean;
    /**
     * 
     * @type {string}
     * @memberof AgentInstanceParameters
     */
    jwtVerifyKeyId: string;
}

/**
 * Check if a given object implements the AgentInstanceParameters interface.
 */
export function instanceOfAgentInstanceParameters(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "endpoint" in value;
    isInstance = isInstance && "buildId" in value;
    isInstance = isInstance && "configVersion" in value;
    isInstance = isInstance && "state" in value;
    isInstance = isInstance && "tlsCertificateId" in value;
    isInstance = isInstance && "tlsCertificateSignedByPublicCa" in value;
    isInstance = isInstance && "jwtVerifyKeyId" in value;

    return isInstance;
}

export function AgentInstanceParametersFromJSON(json: any): AgentInstanceParameters {
    return AgentInstanceParametersFromJSONTyped(json, false);
}

export function AgentInstanceParametersFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentInstanceParameters {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'endpoint': json['endpoint'],
        'buildId': json['buildId'],
        'configVersion': json['configVersion'],
        'state': AgentInstanceStateFromJSON(json['state']),
        'tlsCertificateId': json['tlsCertificateId'],
        'tlsCertificateSignedByPublicCa': json['tlsCertificateSignedByPublicCa'],
        'jwtVerifyKeyId': json['jwtVerifyKeyId'],
    };
}

export function AgentInstanceParametersToJSON(value?: AgentInstanceParameters | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'endpoint': value.endpoint,
        'buildId': value.buildId,
        'configVersion': value.configVersion,
        'state': AgentInstanceStateToJSON(value.state),
        'tlsCertificateId': value.tlsCertificateId,
        'tlsCertificateSignedByPublicCa': value.tlsCertificateSignedByPublicCa,
        'jwtVerifyKeyId': value.jwtVerifyKeyId,
    };
}

