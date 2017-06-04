module Edit exposing (Edit, load, create, revert, map, latest)


type Edit a
    = Saved a
    | Edited (Maybe a) a


load : a -> Edit a
load thing =
    Saved thing


create : a -> Edit a
create thing =
    Edited Nothing thing


revert : Edit a -> Maybe (Edit a)
revert thing =
    case thing of
        Saved _ ->
            Just thing

        Edited original _ ->
            case original of
                Nothing ->
                    Nothing

                Just something ->
                    Just <| Saved something


map : (a -> a) -> Edit a -> Edit a
map change thing =
    case thing of
        Saved saved ->
            Edited (Just saved) (change saved)

        Edited original updated ->
            Edited original (change updated)


latest : Edit a -> a
latest thing =
    case thing of
        Saved saved ->
            saved

        Edited _ updated ->
            updated
