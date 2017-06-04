module Party exposing (..)

import Html exposing (..)
import Html.CssHelpers
import PinCss as PCss
import Dict exposing (Dict)
import Edit exposing (Edit)
import Party.Pin as Pin exposing (..)
import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline as Pipeline
import Json.Encode as Encode exposing (Value)


{ id, class, classList } =
    Html.CssHelpers.withNamespace PCss.pinspace


aClass : a -> Attribute msg
aClass someClass =
    class [ someClass ]


type alias Party =
    { hash : String
    , description : String
    , pins : Dict String (Edit Pin)
    }


init : String -> String -> Party
init hash description =
    Party hash description Dict.empty


type Msg
    = UpdateDescription String
    | LoadPins (List Pin)
    | ReloadPins (List Pin)
    | LoadPin Pin
    | CreatePin Pin
    | RevertPin String
    | UpdatePin String Pin.Msg


update : Msg -> Party -> Party
update msg party =
    case msg of
        UpdateDescription description ->
            { party | description = description }

        LoadPins pins ->
            reloadPins pins party

        ReloadPins pins ->
            reloadPins pins party

        LoadPin pin ->
            replacePin (Edit.load pin) party

        CreatePin pin ->
            replacePin (Edit.create pin) party

        RevertPin pinHash ->
            revertPin pinHash party

        UpdatePin pinHash pinMsg ->
            mapPin (Pin.update pinMsg) pinHash party


view : Party -> Html Msg
view party =
    div
        [ aClass PCss.PartyEntity ]
        [ div [ aClass PCss.PartyDescription ] [ text party.description ]
        , div [ aClass PCss.PartyHash ] [ text party.hash ]
        , div
            [ aClass PCss.PartyPins ]
            (List.map viewPin (currentPins party))
        ]


viewPin : Pin -> Html Msg
viewPin pin =
    Html.map (UpdatePin pin.hash) (Pin.view pin)


pins : Party -> List (Edit Pin)
pins party =
    Dict.values party.pins


currentPins : Party -> List Pin
currentPins party =
    List.map Edit.latest (pins party)


reloadPins : List Pin -> Party -> Party
reloadPins pins party =
    let
        emptyParty =
            { party | pins = Dict.empty }

        hashedPins =
            List.map
                (\p ->
                    ( p.hash, (Edit.load p) )
                )
                pins
    in
        { party | pins = Dict.fromList hashedPins }


replacePin : Edit Pin -> Party -> Party
replacePin pin party =
    let
        pinHash =
            .hash (Edit.latest pin)
    in
        { party
            | pins = Dict.insert pinHash pin party.pins
        }


removePin : String -> Party -> Party
removePin pinHash party =
    { party
        | pins = Dict.remove pinHash party.pins
    }


getEditPin : String -> Party -> Maybe (Edit Pin)
getEditPin pinHash party =
    Dict.get pinHash party.pins


mapPin : (Pin -> Pin) -> String -> Party -> Party
mapPin change pinHash party =
    case (getEditPin pinHash party) of
        Nothing ->
            party

        Just editPin ->
            replacePin (Edit.map change editPin) party


revertPin : String -> Party -> Party
revertPin pinHash party =
    case (getEditPin pinHash party) of
        Nothing ->
            party

        Just pin ->
            let
                reverted =
                    Edit.revert pin
            in
                case reverted of
                    Nothing ->
                        removePin pinHash party

                    Just something ->
                        replacePin something party


decode : Decoder Party
decode =
    Pipeline.decode Party
        |> Pipeline.required "hash" Decode.string
        |> Pipeline.required "description" Decode.string
        |> Pipeline.hardcoded Dict.empty


mutableParts : Party -> List ( String, Value )
mutableParts party =
    [ ( "description", Encode.string party.description )
    ]


encode : Party -> Value
encode party =
    Encode.object (mutableParts party)


encodeCreate : Party -> Value
encodeCreate party =
    Encode.object
        ([ ( "hash", Encode.string party.hash )
         ]
            ++ (mutableParts party)
        )
