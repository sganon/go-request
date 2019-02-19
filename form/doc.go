/*
Package form provides decoding functions for an http request forms (query, url-encoded, multipart).
The decoding consists of both unmarshalling of the query and validation.
Validation consist of builtin tag options like required but also by
self validation of your destination interface.
Customs types can be used for unmarshalling if they implements encoding.TextUnmarshaller, encoding.BinaryUnmarshaler, StringSetter.
*/
package form
