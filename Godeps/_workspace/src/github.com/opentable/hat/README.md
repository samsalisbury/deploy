# hat - hat's an api that's an api t

Idea: Define basic CRUD operations on resources, then define HTTP behaviour in terms of those operations. Then write minimal code to have those resources presented in hypertext formats, such as HAL.

Status: Basic structure nearly there, some features implemented:

- Embedding of simple members, with optional field filters
- Embedding of collections, with pagination & item field filtering
- Linking to logical children

Required features:

- Filtering from query-string parameters
- Richer linking
- Supporting child slices (presently only maps supported)

Dirty code that needs cleaning:

- `type IN int` should be `type IN struct {...}` containing all the info necessary to bind and render inputs. Presently the list of IN_X consts is referenced in multiple places, which is messy and error prone.
- Passing of anything to do with multiple inputs should always be done using slices, because ordering matters. Currently, the way ordering is derived is very hacky.
- Almost everything is exported. This makes it easier to experiment on at this early stage, but once the API is closer to full-featured, expect most types to be made package-visible only.

