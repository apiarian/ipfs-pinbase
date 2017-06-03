module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Html.CssHelpers
import PinCss as PCss
import Http
import Dict exposing (Dict)
import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline as Pipeline
import Json.Encode as Encode exposing (Value)


{ id, class, classList } =
    Html.CssHelpers.withNamespace PCss.pinspace


main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


type alias Model =
    { parties : List Party
    , pins : Dict String (List (Edit Pin))
    }


emptyModel : Model
emptyModel =
    Model [] Dict.empty


type Edit a
    = Original a
    | Edited a a


revertEdit : Edit a -> Edit a
revertEdit thing =
    case thing of
        Original something ->
            Original something

        Edited original new ->
            Original original


makeEdit : (a -> a) -> Edit a -> Edit a
makeEdit f thing =
    case thing of
        Original something ->
            Edited something (f something)

        Edited original old ->
            Edited original (f old)


type alias Party =
    { hash : String
    , description : String
    }


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


type Msg
    = GetParties
    | NewParties (Result Http.Error (List Party))
    | GetPins String
    | NewPins String (Result Http.Error (List (Edit Pin)))
    | EditPin String String
    | FlipWantPinned String String
    | UpdatePinAliases String String String
    | RevertPin String String
    | CommitPin String String
    | UpdatePin String String (Result Http.Error (Edit Pin))


init : ( Model, Cmd Msg )
init =
    ( emptyModel, getParties )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GetParties ->
            ( model, getParties )

        NewParties (Ok parties) ->
            { model
                | parties = parties
            }
                ! List.map (\p -> getPins p.hash) parties

        NewParties (Err _) ->
            ( model, Cmd.none )

        GetPins partyHash ->
            ( model, getPins partyHash )

        NewPins partyHash (Ok pins) ->
            ( { model
                | pins = Dict.insert partyHash pins model.pins
              }
            , Cmd.none
            )

        NewPins _ (Err _) ->
            ( model, Cmd.none )

        FlipWantPinned partyHash pinHash ->
            ( updatePin model partyHash pinHash (\p -> { p | wantPinned = not p.wantPinned })
            , Cmd.none
            )

        UpdatePinAliases partyHash pinHash newAliases ->
            ( updatePin model partyHash pinHash (\p -> { p | aliases = textToAliases newAliases })
            , Cmd.none
            )

        EditPin partyHash pinHash ->
            ( updatePin model partyHash pinHash (\p -> p)
            , Cmd.none
            )

        RevertPin partyHash pinHash ->
            ( revertPin model partyHash pinHash
            , Cmd.none
            )

        CommitPin partyHash pinHash ->
            ( model, commitPin model partyHash pinHash )

        UpdatePin partyHash pinHash (Ok pin) ->
            ( setPin model partyHash pinHash pin, Cmd.none )

        UpdatePin partyHash pinHash (Err _) ->
            ( model, Cmd.none )


setPin : Model -> String -> String -> Edit Pin -> Model
setPin model partyHash pinHash pin =
    { model
        | pins = Dict.update partyHash (updateMaybe (pinSelector pinHash) (\_ -> pin)) model.pins
    }


updatePin : Model -> String -> String -> (Pin -> Pin) -> Model
updatePin model partyHash pinHash change =
    { model
        | pins = Dict.update partyHash (updateMaybe (pinSelector pinHash) (makeEdit change)) model.pins
    }


revertPin : Model -> String -> String -> Model
revertPin model partyHash pinHash =
    { model
        | pins = Dict.update partyHash (updateMaybe (pinSelector pinHash) revertEdit) model.pins
    }


pinSelector : String -> Edit Pin -> Bool
pinSelector pinHash pin =
    let
        hash =
            case pin of
                Original pin ->
                    pin.hash

                Edited original new ->
                    original.hash
    in
        pinHash == hash


updateMaybe : (a -> Bool) -> (a -> a) -> Maybe (List a) -> Maybe (List a)
updateMaybe selector change list =
    case list of
        Nothing ->
            Nothing

        Just list ->
            Just (List.map (conditionalChange selector change) list)


conditionalChange : (a -> Bool) -> (a -> a) -> a -> a
conditionalChange selector change thing =
    if selector thing then
        change thing
    else
        thing


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


aClass : a -> Attribute msg
aClass someClass =
    class [ someClass ]


view : Model -> Html Msg
view model =
    div
        [ class [ PCss.Main ] ]
        [ div
            [ aClass PCss.PartyList ]
            (List.map (viewParty model) model.parties)
        ]


viewParty : Model -> Party -> Html Msg
viewParty model party =
    let
        pins =
            Maybe.withDefault [] (Dict.get party.hash model.pins)
    in
        div
            [ aClass PCss.PartyEntity ]
            [ div [ aClass PCss.PartyDescription ] [ text party.description ]
            , div [ aClass PCss.PartyHash ] [ text party.hash ]
            , div
                [ aClass PCss.PartyPins ]
                (List.map (viewPin party.hash) pins)
            ]


viewPin : String -> Edit Pin -> Html Msg
viewPin partyHash pin =
    let
        ( latestPin, isEditing ) =
            case pin of
                Original pin ->
                    ( pin, False )

                Edited original new ->
                    ( new, True )

        confirmBlock =
            if isEditing then
                div
                    []
                    [ button
                        [ Html.Attributes.class "pure-button"
                        , onClick <| CommitPin partyHash latestPin.hash
                        ]
                        [ text "Commit" ]
                    , button
                        [ Html.Attributes.class "pure-button"
                        , onClick <| RevertPin partyHash latestPin.hash
                        ]
                        [ text "Revert" ]
                    ]
            else
                text ""
    in
        div
            [ classList
                [ ( PCss.PinEntity, True )
                , ( PCss.PinEditing, isEditing )
                ]
            ]
            [ div
                [ aClass PCss.PinHeader ]
                [ div [ aClass PCss.PinName ] [ text <| pinAlias latestPin ]
                , div [ aClass PCss.PinHash ] [ text latestPin.hash ]
                ]
            , viewAliases partyHash latestPin isEditing
            , viewPinStatus partyHash latestPin
            , confirmBlock
            ]


aliasesToText : List String -> String
aliasesToText aliases =
    String.join "\n" aliases


textToAliases : String -> List String
textToAliases string =
    List.filter (not << String.isEmpty) <|
        List.map String.trim <|
            String.split "\n" string


viewAliases : String -> Pin -> Bool -> Html Msg
viewAliases partyHash pin isEditing =
    let
        aliases =
            if isEditing then
                textarea
                    [ rows ((List.length pin.aliases) + 1)
                    , onInput <| UpdatePinAliases partyHash pin.hash
                    ]
                    [ text <| aliasesToText pin.aliases ]
            else
                ul
                    [ aClass PCss.PinAliasesList ]
                    (List.map
                        (\a -> li [ aClass PCss.PinAlias ] [ text a ])
                        pin.aliases
                    )

        editButton =
            if isEditing then
                text ""
            else
                button [ onClick <| EditPin partyHash pin.hash ] [ text "edit" ]
    in
        div
            [ aClass PCss.PinAliasesBox ]
            [ p
                [ aClass PCss.PinAliasesLabel ]
                [ text "Aliases: "
                , editButton
                ]
            , aliases
            ]


viewPinStatus : String -> Pin -> Html Msg
viewPinStatus partyHash pin =
    let
        isPinned =
            pin.status == Pinned

        ( statusClass, statusString ) =
            case pin.status of
                Pinned ->
                    ( PCss.PinStatusGood, "ok" )

                Pending ->
                    ( PCss.PinStatusPending, "pending" )

                Unpinned ->
                    ( PCss.PinStatusGood, "ok" )

                Error ->
                    ( PCss.PinStatusError, "error: " ++ pin.lastError )

                Fatal ->
                    ( PCss.PinStatusError, "fatal error: " ++ pin.lastError )
    in
        div
            []
            [ div
                [ aClass PCss.PinPinBox ]
                [ label
                    []
                    [ text "pin"
                    , input
                        [ type_ "checkbox"
                        , checked pin.wantPinned
                        , onClick (FlipWantPinned partyHash pin.hash)
                        ]
                        []
                    ]
                ]
            , div
                [ class [ PCss.PinStatusText, statusClass ] ]
                [ text statusString ]
            ]


pinAlias : Pin -> String
pinAlias pin =
    let
        mainAlias =
            List.head pin.aliases
    in
        case mainAlias of
            Just something ->
                something

            Nothing ->
                "-no primary alias-"


getLatestPin : Model -> String -> String -> Maybe Pin
getLatestPin model partyHash pinHash =
    let
        pins =
            Dict.get partyHash model.pins
    in
        case pins of
            Nothing ->
                Nothing

            Just pins ->
                let
                    filteredPins =
                        List.filter (pinSelector pinHash) pins
                in
                    case filteredPins of
                        [] ->
                            Nothing

                        x :: _ ->
                            case x of
                                Original pin ->
                                    Just pin

                                Edited _ pin ->
                                    Just pin


commitPin : Model -> String -> String -> Cmd Msg
commitPin model partyHash pinHash =
    let
        pin =
            getLatestPin model partyHash pinHash
    in
        case pin of
            Nothing ->
                Cmd.none

            Just pin ->
                Http.send (UpdatePin partyHash pinHash)
                    (Http.request
                        { method = "PATCH"
                        , headers =
                            [ Http.header "content-type" "application/json"
                            ]
                        , url = pinUrl partyHash pinHash
                        , body = Http.jsonBody <| encodePin pin
                        , expect = Http.expectJson <| editDecoder decodePin
                        , timeout = Nothing
                        , withCredentials = False
                        }
                    )


encodePin : Pin -> Value
encodePin pin =
    Encode.object
        [ ( "aliases", Encode.list (List.map Encode.string pin.aliases) )
        , ( "want-pinned", Encode.bool pin.wantPinned )
        ]


getParties : Cmd Msg
getParties =
    Http.send NewParties (Http.get partiesUrl decodeParties)


apiBaseUrl : String
apiBaseUrl =
    "http://localhost:3000/api/"


partiesUrl : String
partiesUrl =
    apiBaseUrl ++ "parties"


pinsUrl : String -> String
pinsUrl partyHash =
    partiesUrl ++ "/" ++ partyHash ++ "/pins"


pinUrl : String -> String -> String
pinUrl partyHash pinHash =
    (pinsUrl partyHash) ++ "/" ++ pinHash


getPins : String -> Cmd Msg
getPins partyHash =
    Http.send (NewPins partyHash) (Http.get (pinsUrl partyHash) decodePins)


decodeParties : Decoder (List Party)
decodeParties =
    Decode.list decodeParty


decodeParty : Decoder Party
decodeParty =
    Pipeline.decode Party
        |> Pipeline.required "hash" Decode.string
        |> Pipeline.required "description" Decode.string


decodePins : Decoder (List (Edit Pin))
decodePins =
    Decode.list (editDecoder decodePin)


editDecoder : Decoder a -> Decoder (Edit a)
editDecoder decoder =
    Decode.map Original decoder


decodePin : Decoder Pin
decodePin =
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
