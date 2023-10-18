/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.1
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { AgentConfigurationParameters } from './AgentConfigurationParameters';
import {
    AgentConfigurationParametersFromJSON,
    AgentConfigurationParametersFromJSONTyped,
    AgentConfigurationParametersToJSON,
} from './AgentConfigurationParameters';

/**
 * 
 * @export
 * @interface AgentConfiguration
 */
export interface AgentConfiguration {
    /**
     * 
     * @type {AgentConfigurationParameters}
     * @memberof AgentConfiguration
     */
    config: AgentConfigurationParameters;
    /**
     * Version of the agent configuration
     * @type {string}
     * @memberof AgentConfiguration
     */
    version: string;
    /**
     * 
     * @type {string}
     * @memberof AgentConfiguration
     */
    nextRefreshToken?: string;
    /**
     * 
     * @type {Date}
     * @memberof AgentConfiguration
     */
    nextRefreshAfter?: Date;
}

/**
 * Check if a given object implements the AgentConfiguration interface.
 */
export function instanceOfAgentConfiguration(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "config" in value;
    isInstance = isInstance && "version" in value;

    return isInstance;
}

export function AgentConfigurationFromJSON(json: any): AgentConfiguration {
    return AgentConfigurationFromJSONTyped(json, false);
}

export function AgentConfigurationFromJSONTyped(json: any, ignoreDiscriminator: boolean): AgentConfiguration {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'config': AgentConfigurationParametersFromJSON(json['config']),
        'version': json['version'],
        'nextRefreshToken': !exists(json, 'nextRefreshToken') ? undefined : json['nextRefreshToken'],
        'nextRefreshAfter': !exists(json, 'nextRefreshAfter') ? undefined : (new Date(json['nextRefreshAfter'])),
    };
}

export function AgentConfigurationToJSON(value?: AgentConfiguration | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'config': AgentConfigurationParametersToJSON(value.config),
        'version': value.version,
        'nextRefreshToken': value.nextRefreshToken,
        'nextRefreshAfter': value.nextRefreshAfter === undefined ? undefined : (value.nextRefreshAfter.toISOString()),
    };
}

