module Party.Pin exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Html.CssHelpers
import PinCss as PCss
import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline as Pipeline
import Json.Encode as Encode exposing (Value)


{ id, class, classList } =
    Html.CssHelpers.withNamespace PCss.pinspace


aClass : a -> Attribute msg
aClass someClass =
    class [ someClass ]


type alias Pin =
    { hash : String
    , aliases : List String
    , wantPinned : Bool
    , status : PinStatus
    , lastError : String
    }


type PinStatus
    = Pending
    | Pinned
    | Unpinned
    | Error
    | Fatal
    | Unknown


init : String -> String -> Bool -> Pin
init hash primaryAlias wantPinned =
    Pin hash [ primaryAlias ] wantPinned Unknown ""


type Msg
    = ToggleWantPinned
    | ReplaceAliases (List String)
    | UpdateAliases String


update : Msg -> Pin -> Pin
update msg pin =
    case msg of
        ToggleWantPinned ->
            { pin | wantPinned = not pin.wantPinned }

        ReplaceAliases new ->
            { pin | aliases = new }

        UpdateAliases text ->
            { pin | aliases = textToList text }


view : Pin -> Html Msg
view pin =
    div
        [ aClass PCss.PinEntity ]
        [ div
            [ aClass PCss.PinPinBox ]
            [ label
                []
                [ text "pin"
                , input
                    [ type_ "checkbox"
                    , checked pin.wantPinned
                    , onClick ToggleWantPinned
                    ]
                    []
                ]
            ]
        ]


aliasesText : Pin -> String
aliasesText pin =
    listToText pin.aliases


textToList : String -> List String
textToList text =
    List.filter (not << String.isEmpty) <|
        List.map String.trim <|
            String.split "\n" text


listToText : List String -> String
listToText list =
    String.join "\n" list


decode : Decoder Pin
decode =
    Pipeline.decode Pin
        |> Pipeline.required "hash" Decode.string
        |> Pipeline.required "aliases" (Decode.list Decode.string)
        |> Pipeline.required "want-pinned" Decode.bool
        |> Pipeline.required "status" decodePinStatus
        |> Pipeline.required "last-error" Decode.string


decodePinStatus : Decoder PinStatus
decodePinStatus =
    Decode.string |> Decode.andThen stringToPinStatus


stringToPinStatus : String -> Decoder PinStatus
stringToPinStatus s =
    let
        ps =
            case s of
                "pending" ->
                    Just Pending

                "pinned" ->
                    Just Pinned

                "unpinned" ->
                    Just Unpinned

                "error" ->
                    Just Error

                "fatal" ->
                    Just Fatal

                _ ->
                    Nothing
    in
        case ps of
            Just status ->
                Decode.succeed status

            Nothing ->
                Decode.fail <| "unknown status value '" ++ s ++ "'"


mutableParts : Pin -> List ( String, Value )
mutableParts pin =
    [ ( "aliases", Encode.list (List.map Encode.string pin.aliases) )
    , ( "want-pinned", Encode.bool pin.wantPinned )
    ]


encode : Pin -> Value
encode pin =
    Encode.object (mutableParts pin)


encodeCreate : Pin -> Value
encodeCreate pin =
    Encode.object
        ([ ( "hash", Encode.string pin.hash )
         ]
            ++ (mutableParts pin)
        )
