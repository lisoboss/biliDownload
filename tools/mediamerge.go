package tools

import (
	"biliDownload/tools/media"
	"io"
	"os"
)

func MediaMerge(videoPath, audioPath, savePath string) (err error) {
	rVideo, err := media.NewRoot(videoPath)
	if err != nil {
		return
	}
	rAudio, err := media.NewRoot(audioPath)
	if err != nil {
		return
	}

	// TrackExtends
	rvTrackExtends := rVideo.GetBox("trex")
	raTrackExtends := rAudio.GetBox("trex")

	// rVideo add rAudio 'trex'
	raTrackExtendsLength := raTrackExtends.Length()
	_pBox := rvTrackExtends.ParentBox()
	_pBox.AddContainerBox(raTrackExtends)
	_pBox.AddLength(raTrackExtendsLength)

	// Track
	rvTrack := rVideo.GetBox("trak")
	raTrack := rAudio.GetBox("trak")

	// rVideo add rAudio 'trak'
	raTrackLength := raTrack.Length()
	_pBox = rvTrack.ParentBox()
	_pBox.AddContainerBox(raTrack)
	_pBox.AddLength(raTrackLength)

	// new media
	newMedia := new(media.Root)

	// add box
	rvIndex := 0
	rvBoxList := rVideo.BoxList()
	rvBoxListLen := len(rvBoxList)
	for i := 0; i < rvBoxListLen; i++ {
		box := rvBoxList[i]
		if box.Type() == "moof" {
			rvIndex = i
			break
		}
		newMedia.AddBox(box)
	}

	// raIndex
	raIndex := 0
	raBoxList := rAudio.BoxList()
	raBoxListLen := len(raBoxList)
	for i := 0; i < raBoxListLen; i++ {
		box := raBoxList[i]
		if box.Type() == "moof" {
			raIndex = i
			break
		}
	}

	// add rAudio audio data(moof+mdat)
	var sequenceNumber uint32 = 1
	for rvIndex < rvBoxListLen || raIndex < raBoxListLen {
		if rvIndex < rvBoxListLen {
			box := rvBoxList[rvIndex]
			updateMovieFragmentHeaderBox(box, sequenceNumber)
			sequenceNumber++
			newMedia.AddBox(box)
			rvIndex++
			newMedia.AddBox(rvBoxList[rvIndex])
			rvIndex++
		}
		if raIndex < raBoxListLen {
			box := raBoxList[raIndex]
			updateMovieFragmentHeaderBox(box, sequenceNumber)
			sequenceNumber++
			newMedia.AddBox(box)
			raIndex++
			newMedia.AddBox(raBoxList[raIndex])
			raIndex++
		}
	}

	// save newMedia
	fileOut, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer func() {
		_ = fileOut.Close()
	}()
	_, err = fileOut.Write(newMedia.Dump())
	return
}

func MediaMergeFromReader(videoReader, audioReader *Reader, savePath string) (err error) {
	// create dir
	if err = CreateDirFromFilePath(savePath); err != nil {
		Log.Fatal(err)
	}
	// save
	fileOut, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer func() {
		_ = fileOut.Close()
	}()

	// get 'moov'
	var (
		vMoov *media.Box
		aMoov *media.Box
	)
	for true {
		box, vErr := media.LoadBoxFrom(videoReader)
		if vErr != nil {
			return vErr
		}
		if box.Type() == "moov" {
			vMoov = box
			break
		}
		// save box
		if _, err = fileOut.Write(box.Dump()); err != nil {
			return
		}
	}
	for true {
		box, vErr := media.LoadBoxFrom(audioReader)
		if vErr != nil {
			return vErr
		}
		if box.Type() == "moov" {
			aMoov = box
			break
		}
	}

	// TrackExtends
	rvTrackExtends := vMoov.GetBox("trex")
	raTrackExtends := aMoov.GetBox("trex")

	// rVideo add rAudio 'trex'
	raTrackExtendsLength := raTrackExtends.Length()
	_pBox := rvTrackExtends.ParentBox()
	_pBox.AddContainerBox(raTrackExtends)
	_pBox.AddLength(raTrackExtendsLength)

	// Track
	rvTrack := vMoov.GetBox("trak")
	raTrack := aMoov.GetBox("trak")

	// rVideo add rAudio 'trak'
	raTrackLength := raTrack.Length()
	_pBox = rvTrack.ParentBox()
	_pBox.AddContainerBox(raTrack)
	_pBox.AddLength(raTrackLength)

	var sequenceNumber uint32 = 1

	// save vMoov box
	if _, err = fileOut.Write(vMoov.Dump()); err != nil {
		return
	}

	// to moof
	for true {
		box, vErr := media.LoadBoxFrom(videoReader)
		if vErr != nil {
			return vErr
		}
		if box.Type() == "moof" {
			// save moof
			updateMovieFragmentHeaderBox(box, sequenceNumber)
			sequenceNumber++
			if _, err = fileOut.Write(box.Dump()); err != nil {
				return
			}

			// save mdat
			box, err = media.LoadBoxFrom(videoReader)
			if err != nil {
				return err
			}
			if _, err = fileOut.Write(box.Dump()); err != nil {
				return
			}
			break
		}
		// save box
		if _, err = fileOut.Write(box.Dump()); err != nil {
			return
		}
	}
	for true {
		box, vErr := media.LoadBoxFrom(audioReader)
		if vErr != nil {
			return vErr
		}
		if box.Type() == "moof" {
			// save moof
			updateMovieFragmentHeaderBox(box, sequenceNumber)
			sequenceNumber++
			if _, err = fileOut.Write(box.Dump()); err != nil {
				return
			}

			// save mdat
			box, err = media.LoadBoxFrom(audioReader)
			if err != nil {
				return err
			}
			if _, err = fileOut.Write(box.Dump()); err != nil {
				return
			}
			break
		}
	}

	// save data(moof+mdat)
	for !videoReader.Over() || !audioReader.Over() {
		if !videoReader.Over() {
			if vBox, vErr := media.LoadBoxFrom(videoReader); vErr != nil {
				if vErr != io.EOF {
					return vErr
				}
			} else {
				// save moof
				updateMovieFragmentHeaderBox(vBox, sequenceNumber)
				sequenceNumber++
				if _, err = fileOut.Write(vBox.Dump()); err != nil {
					return
				}

				// save mdat
				vBox, err = media.LoadBoxFrom(videoReader)
				if err != nil {
					return err
				}
				if _, err = fileOut.Write(vBox.Dump()); err != nil {
					return
				}
			}
		}
		if !audioReader.Over() {
			if vBox, vErr := media.LoadBoxFrom(audioReader); vErr != nil {
				if vErr != io.EOF {
					return vErr
				}
			} else {
				// save moof
				updateMovieFragmentHeaderBox(vBox, sequenceNumber)
				sequenceNumber++
				if _, err = fileOut.Write(vBox.Dump()); err != nil {
					return
				}

				// save mdat
				vBox, err = media.LoadBoxFrom(audioReader)
				if err != nil {
					return err
				}
				if _, err = fileOut.Write(vBox.Dump()); err != nil {
					return
				}
			}
		}
	}

	return
}

func updateMovieFragmentHeaderBox(box *media.Box, sequenceNumber uint32) {
	if box == nil {
		return
	}
	_mbox := box.GetBox("mfhd")
	if _mbox == nil {
		return
	}
	_movieFragmentHeaderBox := new(media.MovieFragmentHeaderBox)
	_movieFragmentHeaderBox.Load(_mbox.Dump())
	_movieFragmentHeaderBox.SetSequenceNumber(sequenceNumber)
	_mbox.Load(_movieFragmentHeaderBox.Dump())
}
