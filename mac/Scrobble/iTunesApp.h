//
//  iTunesApp.h
//  Scrobble
//
//  Created by Nishanth Shanmugham on 5/8/18.
//  Copyright © 2018 Nishanth Shanmugham. All rights reserved.
//

#ifndef iTunesApp_h
#define iTunesApp_h

#import "iTunes.h" // generated via: sdef /Applications/iTunes.app | sdp -fh --basename iTunes

@interface iTunesApp : NSObject

@property iTunesApplication *app;
- (BOOL) currentTrackIsSong;

@end

#endif /* iTunesApp_h */
